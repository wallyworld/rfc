// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tls_test

import (
	stdtls "crypto/tls"
	"crypto/x509"

	"github.com/juju/errors"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/rfc/tls"
	"github.com/juju/rfc/tls/tlstest"
)

type ConfigSuite struct {
	testing.IsolationSuite
}

var _ = gc.Suite(ConfigSuite{})

func (ConfigSuite) TestTLSFull(c *gc.C) {
	cfg := tls.Config{
		RawCert: tls.RawCert{
			CertPEM:   tlstest.ServerCert.CertPEM,
			KeyPEM:    tlstest.ServerCert.KeyPEM,
			CACertPEM: tlstest.CACert.CertPEM,
		},
		ServerName:            "a.b.c",
		ExpectedServerCertPEM: validCert2,
	}
	cert, err := cfg.Cert()
	c.Assert(err, jc.ErrorIsNil)
	serverCAs := x509.NewCertPool()
	caCert, err := cfg.CACert()
	c.Assert(err, jc.ErrorIsNil)
	serverCAs.AddCert(caCert)

	tlsConfig, err := cfg.TLS()
	c.Assert(err, jc.ErrorIsNil)

	c.Check(tlsConfig, jc.DeepEquals, &stdtls.Config{
		CipherSuites: secureConfig.CipherSuites,
		MinVersion:   secureConfig.MinVersion,
		ServerName:   "a.b.c",
		RootCAs:      serverCAs,
		Certificates: []stdtls.Certificate{
			cert,
		},
		NameToCertificate: map[string]*stdtls.Certificate{
			"*": &cert,
		},
	})
}

func (s *ConfigSuite) TestRawValidateFull(c *gc.C) {
	cfg := tls.Config{
		RawCert: tls.RawCert{
			CertPEM:   tlstest.ServerCert.CertPEM,
			KeyPEM:    tlstest.ServerCert.KeyPEM,
			CACertPEM: tlstest.CACert.CertPEM,
		},
		ServerName:            "a.b.c",
		ExpectedServerCertPEM: validCert2,
	}

	err := cfg.Validate()

	c.Check(err, jc.ErrorIsNil)
}

func (s *ConfigSuite) TestRawValidateNoServerName(c *gc.C) {
	cfg := tls.Config{
		RawCert: tls.RawCert{
			CertPEM:   tlstest.ServerCert.CertPEM,
			KeyPEM:    tlstest.ServerCert.KeyPEM,
			CACertPEM: tlstest.CACert.CertPEM,
		},
		ExpectedServerCertPEM: validCert2,
	}

	err := cfg.Validate()

	c.Check(err, jc.ErrorIsNil)
}

func (s *ConfigSuite) TestRawValidateBadRawCert(c *gc.C) {
	cfg := tls.Config{
		RawCert: tls.RawCert{
			CertPEM:   invalidCert,
			KeyPEM:    validKey,
			CACertPEM: validCACert,
		},
		ServerName:            "a.b.c",
		ExpectedServerCertPEM: validCert2,
	}

	err := cfg.Validate()

	c.Check(err, jc.Satisfies, errors.IsNotValid)
}

func (s *ConfigSuite) TestRawValidateMissingServerCert(c *gc.C) {
	cfg := tls.Config{
		RawCert: tls.RawCert{
			CertPEM:   validCert,
			KeyPEM:    validKey,
			CACertPEM: validCACert,
		},
		ServerName:            "a.b.c",
		ExpectedServerCertPEM: "",
	}

	err := cfg.Validate()

	c.Check(err, jc.Satisfies, errors.IsNotValid)
	c.Check(err, gc.ErrorMatches, `empty ExpectedServerCertPEM`)
}

func (s *ConfigSuite) TestRawValidateBadServerCert(c *gc.C) {
	cfg := tls.Config{
		RawCert: tls.RawCert{
			CertPEM:   validCert,
			KeyPEM:    validKey,
			CACertPEM: validCACert,
		},
		ServerName:            "a.b.c",
		ExpectedServerCertPEM: invalidCert,
	}

	err := cfg.Validate()

	c.Check(err, jc.Satisfies, errors.IsNotValid)
	c.Check(err, gc.ErrorMatches, `invalid ExpectedServerCertPEM: asn1: syntax error: data truncated`)
}

func (s *ConfigSuite) TestRawValidateBadServerCertFormat(c *gc.C) {
	cfg := tls.Config{
		RawCert: tls.RawCert{
			CertPEM:   validCert,
			KeyPEM:    validKey,
			CACertPEM: validCACert,
		},
		ServerName:            "a.b.c",
		ExpectedServerCertPEM: "abc",
	}

	err := cfg.Validate()

	c.Check(err, jc.Satisfies, errors.IsNotValid)
	c.Check(err, gc.ErrorMatches, `invalid ExpectedServerCertPEM: no certificates found`)
}

var secureConfig = tls.SecureTLSConfig()
