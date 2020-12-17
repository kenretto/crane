package proxy

import (
	"context"
	"golang.org/x/crypto/ssh"
	"net"
	"net/http"
	"time"
)

// SSHProxy ssh 代理
type SSHProxy struct {
	Network, Address string
	*ssh.ClientConfig
	client *ssh.Client

	dialed bool
}

// NewSSHProxy 初始化 ssh 代理
func NewSSHProxy(network, address string) *SSHProxy {
	return &SSHProxy{
		Network: network,
		Address: address,
		ClientConfig: &ssh.ClientConfig{
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			HostKeyAlgorithms: []string{
				ssh.KeyAlgoRSA,
				ssh.KeyAlgoDSA,
				ssh.KeyAlgoECDSA256,
				ssh.KeyAlgoECDSA384,
				ssh.KeyAlgoECDSA521,
				ssh.KeyAlgoED25519,
			},
		},
	}
}

// AuthByUsernameAndPassword 基于用户名和密码的ssh登录认证
func (sshProxy *SSHProxy) AuthByUsernameAndPassword(username, password string) *SSHProxy {
	sshProxy.User = username
	sshProxy.Auth = append(sshProxy.Auth, ssh.Password(password))
	return sshProxy
}

// SetHostKeyCallback set host key callback
func (sshProxy *SSHProxy) SetHostKeyCallback(callback ssh.HostKeyCallback) *SSHProxy {
	sshProxy.HostKeyCallback = callback
	return sshProxy
}

// SetHostKeyAlgorithms set host key algorithms
func (sshProxy *SSHProxy) SetHostKeyAlgorithms(HostKeyAlgorithms []string) *SSHProxy {
	sshProxy.HostKeyAlgorithms = HostKeyAlgorithms
	return sshProxy
}

// SetTimeout set timeout
func (sshProxy *SSHProxy) SetTimeout(duration time.Duration) *SSHProxy {
	sshProxy.Timeout = duration
	return sshProxy
}

// Dial 连接
func (sshProxy *SSHProxy) Dial() (err error) {
	sshProxy.client, err = ssh.Dial(sshProxy.Network, sshProxy.Address, sshProxy.ClientConfig)
	if err == nil {
		sshProxy.dialed = true
	}
	return err
}

// HTTPTransport 返回 *http.Transport, 如果没有进行过连接, 将会连接, dial 错误时将 panic
func (sshProxy *SSHProxy) HTTPTransport() *http.Transport {
	if !sshProxy.dialed {
		err := sshProxy.Dial()
		if err != nil {
			panic(err)
		}
	}
	var transport = new(http.Transport)
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return sshProxy.client.Dial(network, addr)
	}
	return transport
}
