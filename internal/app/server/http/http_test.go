package http_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	server "github.com/patraden/ya-practicum-go-shortly/internal/app/server/http"
)

func testRouter(t *testing.T) http.Handler {
	t.Helper()

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte("pong"))
		assert.NoError(t, err)
	})

	return mux
}

func TestServerRunAndShutdown(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		ServerAddr:              ":8081", // Use a different port for testing
		ServerReadHeaderTimeout: 5 * time.Second,
		ServerWriteTimeout:      10 * time.Second,
		ServerIdleTimeout:       15 * time.Second,
		EnableHTTPS:             false,
	}

	handler := http.NewServeMux()
	server := server.NewServer(cfg, handler)
	ctx := context.Background()
	wgr := sync.WaitGroup{}

	// Run the server in a separate goroutine

	wgr.Add(1)

	go func() {
		defer wgr.Done()

		err := server.Run()
		assert.ErrorIs(t, err, http.ErrServerClosed)
	}()

	// Wait a bit to ensure the server has started
	time.Sleep(500 * time.Millisecond)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+cfg.ServerAddr, nil)
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)

	require.NoError(t, err)

	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	require.NoError(t, err)

	wgr.Wait()
}

func generateTempTLSFiles(t *testing.T) (string, string) {
	t.Helper()

	certFile, err := os.CreateTemp("", "test_cert_*.pem")
	require.NoError(t, err)

	keyFile, err := os.CreateTemp("", "test_key_*.pem")
	require.NoError(t, err)

	t.Cleanup(func() {
		os.Remove(certFile.Name())
		os.Remove(keyFile.Name())
	})

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			Organization: []string{"PatrakhinDenis"},
			Country:      []string{"RS"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(24 * time.Hour),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	// Generate RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err)

	// Create test self-signed certificate
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	require.NoError(t, err)

	// Encode and write certificate
	err = pem.Encode(certFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	require.NoError(t, err)

	// Encode and write private key
	err = pem.Encode(keyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	require.NoError(t, err)

	certFile.Close()
	keyFile.Close()

	return certFile.Name(), keyFile.Name()
}

func loadTLSConfig(t *testing.T, certPath string) *tls.Config {
	t.Helper()

	certData, err := os.ReadFile(certPath)
	require.NoError(t, err)

	// Create a new certificate pool and append our test certificate
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(certData) {
		t.Fatal("failed to append certificate")
		return nil
	}

	return &tls.Config{
		RootCAs:    certPool,
		MinVersion: tls.VersionTLS12,
	}
}

func TestServerRunAndShutdownWithHTTPS(t *testing.T) {
	t.Parallel()

	// Generate temporary TLS cert and key
	certPath, keyPath := generateTempTLSFiles(t)

	cfg := &config.Config{
		ServerAddr:              "127.0.0.1:8443",
		ServerReadHeaderTimeout: 5 * time.Second,
		ServerWriteTimeout:      10 * time.Second,
		ServerIdleTimeout:       15 * time.Second,
		EnableHTTPS:             true,
		TLSCertPath:             certPath,
		TLSKeyPath:              keyPath,
	}

	handler := testRouter(t)
	server := server.NewServer(cfg, handler)
	ctx := context.Background()
	wgr := sync.WaitGroup{}

	wgr.Add(1)

	go func() {
		defer wgr.Done()

		err := server.Run()
		assert.ErrorIs(t, err, http.ErrServerClosed)
	}()

	// Give the server some time to start
	time.Sleep(100 * time.Millisecond)

	tlsConfig := loadTLSConfig(t, certPath)
	require.NotNil(t, tlsConfig)

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+cfg.ServerAddr+"/ping", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = server.Shutdown(shutdownCtx)
	require.NoError(t, err)

	wgr.Wait()
}
