package tunnel

import (
	"errors"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
)

type Forward struct {
	LocalListenerAddr string `yaml:"local_addr"`
	TargetAddr        string `yaml:"target_addr"`
	TunnelAddr        string `yaml:"tunnel_addr"`
	TunnelUser        string `yaml:"tunnel_user"`
	TunnelPwd         string `yaml:"tunnel_pwd"`
	StopCh            chan struct{}
}

func (t Forward) buildTunnelConfig() *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: t.TunnelUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(t.TunnelPwd), // SSH密码
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

func (t Forward) buildTunnelHost() (*ssh.Client, error) {
	return ssh.Dial("tcp", t.TunnelAddr, t.buildTunnelConfig())
}

func (t Forward) buildLocalListener() (net.Listener, error) {
	return net.Listen("tcp", t.LocalListenerAddr)
}

func (t Forward) Run() error {
	localListener, err := t.buildLocalListener()
	if err != nil {
		return err
	}
	log.Printf("initialize local listener success :%s", t.LocalListenerAddr)
	tunnelHost, err := t.buildTunnelHost()
	if err != nil {
		return err
	}
	log.Printf("initialize local tunnel success :%s", t.TunnelAddr)
	go func() {
		<-t.StopCh
		log.Printf("sinal tunnel exit.")
		_ = localListener.Close()
		_ = tunnelHost.Close()
	}()
	for {
		localConn, err := localListener.Accept()
		if errors.Is(err, net.ErrClosed) {
			break
		}
		if err != nil {
			log.Printf("ERROR: Failed to accept local connection: %v", err)
			return err
		}
		remoteConn, err := tunnelHost.Dial("tcp", t.TargetAddr)
		if err != nil {
			log.Printf("ERROR: Failed to dial destination host: %v", err)
			return err
		}
		go func() {
			size, err := io.Copy(localConn, remoteConn)
			if err != nil {
				log.Printf("Error: Failed to copy data: %v", err)
			} else {
				log.Printf("local to target %s ---> %s ---> %s , size: %d", t.LocalListenerAddr, t.TunnelAddr, t.TunnelAddr, size)
			}
		}()
		go func() {
			size, err := io.Copy(remoteConn, localConn)
			if err != nil {
				log.Printf("ERROR: Failed to copy data: %v", err)
			} else {
				log.Printf("target to local %s ---> %s ---> %s, size: %d", t.TargetAddr, t.TunnelAddr, t.LocalListenerAddr, size)
			}
		}()
	}
	return nil
}
