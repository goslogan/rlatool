/*
Copyright © 2024 Nic Gibson <nic.gibson@redis.com>
*/
package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRladmin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rladmin Suite")
}
