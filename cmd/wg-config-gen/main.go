package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path"
	"time"

	"github.com/gitirabassi/vault-plugin-secrets-wireguard/cidrset"
	wgquick "github.com/nmiculinic/wg-quick-go"
	"github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var (
	interfaceName       string
	serverPort          int
	serverAddres        string
	cidr                string
	destinationFolder   string
	clientsAllowedIPs   string
	outputInterface     string
	clientNumber        int
	dnsGateway          bool
	withPresharedKey    bool
	persistentKeepalive string
	logger              = logrus.New()
)

// Config holds the all configurations in memory before being dumped
type Config struct {
	Server  wgquick.Config
	Clients []wgquick.Config
}

func init() {
	flag.StringVar(&interfaceName, "interface-name", "wg0", "wireguard interface name on the server")
	flag.StringVar(&cidr, "cidr", "10.29.100.0/24", "wireguard internal CIDR")
	flag.IntVar(&serverPort, "port", 51820, "port to be exposed by the server")
	flag.StringVar(&persistentKeepalive, "persisten-keepalive", "25s", "Persiste keep alive used by client to refresh connection")
	flag.StringVar(&outputInterface, "output-interface", "eth0", "interface to route all your traffic thru")
	flag.StringVar(&serverAddres, "address", "", "Server public address")
	flag.StringVar(&destinationFolder, "output-dir", "", "Output directory with all configurations (default: tmpdir)")
	flag.StringVar(&clientsAllowedIPs, "clients-allowed-ips", "0.0.0.0/0", "default routes to be used with all clients")
	flag.IntVar(&clientNumber, "num-clients", 1, "number of clients to generate")
	flag.BoolVar(&dnsGateway, "with-dns-gateway", false, "configures gateway IP as DNS server for clients")
	flag.BoolVar(&withPresharedKey, "with-pks", true, "configures server and clients with Preshared Key")
}

func main() {
	flag.Parse()
	if serverAddres == "" {
		logger.Fatalln("'-address' must be defined")
	}
	serverExternalIP := net.ParseIP(serverAddres)
	if serverExternalIP == nil {
		logger.Fatalln("'-address' must be in a form of IPv4 ip")
	}
	pka, err := time.ParseDuration(persistentKeepalive)
	if err != nil {
		logger.Fatalln(err)
	}
	// verify that outputdir is empty or generate tmpdir
	if destinationFolder != "" {
		ok, err := isDirEmpty(destinationFolder)
		if err != nil {
			logger.Fatalln(err)
		}
		if !ok {
			logger.Fatalf("directory '%s' is not empty", destinationFolder)
		}
	} else {
		destinationFolder, err := ioutil.TempDir("tmp", "wg")
		if err != nil {
			logger.Fatal(err)
		}
		logger.Infoln("Created tmpdir at:", destinationFolder)
	}
	_, internalNet, err := net.ParseCIDR(cidr)
	if err != nil {
		logger.Fatalln(err)
	}
	subnet, err := cidrset.NewCIDRSet(internalNet, 32)
	if err != nil {
		logger.Fatalln(err)
	}
	// Skipping first IP as it's usually of type `x.x.x.0`
	_, err = subnet.AllocateNext()
	if err != nil {
		logger.Fatalln(err)
	}
	// First we set all the interface values for the server
	conf := &Config{}
	conf.Server.PostUp = fmt.Sprintf("sysctl -w net.ipv4.ip_forward=1; iptables -A FORWARD -p tcp --dport 22 -i %v -j ACCEPT; iptables -A FORWARD -p tcp --dport 80 -i %v -j ACCEPT; iptables -A FORWARD -p tcp --dport 443 -i %v -j ACCEPT; iptables -t nat -A POSTROUTING -o %v -j MASQUERADE", interfaceName, outputInterface)
	conf.Server.PostDown = fmt.Sprintf("iptables -D FORWARD -i %v -j ACCEPT; iptables -t nat -D POSTROUTING -o %v -j MASQUERADE", interfaceName, outputInterface)
	serverIP, err := subnet.AllocateNext()
	if err != nil {
		logger.Fatalln(err)
	}
	conf.Server.Address = []net.IPNet{*serverIP}
	conf.Server.ListenPort = &serverPort
	serverKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		logger.Fatalln(err)
	}
	conf.Server.PrivateKey = &serverKey

	_, allowedNet, err := net.ParseCIDR(clientsAllowedIPs)
	if err != nil {
		logger.Fatalln(err)
	}
	// Now we configure all the right variables for the clients
	for i := 0; i < clientNumber; i++ {
		client := wgquick.Config{}
		clientKey, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			logger.Fatalln(err)
		}
		client.PrivateKey = &clientKey
		clientIP, err := subnet.AllocateNext()
		if err != nil {
			logger.Fatalln(err)
		}
		client.Address = []net.IPNet{*clientIP}
		if dnsGateway {
			client.DNS = []net.IP{serverIP.IP}
		}

		serverPeer := wgtypes.PeerConfig{
			PublicKey: serverKey.PublicKey(),
			Endpoint: &net.UDPAddr{
				IP:   serverExternalIP,
				Port: serverPort,
			},
			PersistentKeepaliveInterval: &pka,
			AllowedIPs:                  []net.IPNet{*allowedNet},
		}
		if withPresharedKey {
			pskKey, err := wgtypes.GenerateKey()
			if err != nil {
				logger.Fatalln(err)
			}
			serverPeer.PresharedKey = &pskKey
		}
		client.Peers = append(client.Peers, serverPeer)
		conf.Clients = append(conf.Clients, client)
	}
	// now we need to go back to the server to configure the right client peers
	for _, client := range conf.Clients {
		clientPeer := wgtypes.PeerConfig{
			PublicKey:  client.PrivateKey.PublicKey(),
			AllowedIPs: client.Address,
		}
		if withPresharedKey {
			clientPeer.PresharedKey = client.Peers[0].PresharedKey
		}
		conf.Server.Peers = append(conf.Server.Peers, clientPeer)
	}
	// render all config on disk
	err = render(destinationFolder, interfaceName, conf.Server)
	if err != nil {
		logger.Fatalln("couldn't render server conf", err)
	}
	for i, client := range conf.Clients {
		err = render(destinationFolder, fmt.Sprintf("client%v", i), client)
		if err != nil {
			logger.Fatalln("couldn't render client %v configuration", i, err)
		}
	}
	err = renderCorefile(destinationFolder, serverIP.IP.String())
	if err != nil {
		logger.Fatalln("couldn't render server conf", err)
	}
}

func render(destinationFolder, name string, conf wgquick.Config) error {
	fullPath := path.Join(destinationFolder, fmt.Sprintf("%s.conf", name))
	body, err := conf.MarshalText()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fullPath, body, 0644)
}

func isDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
