package infra

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/GlediLami/go-gardener-lab/pkg/cloud"
)

var errBoom = errors.New("cloud is on fire")

func TestEnsureNetwork_CreatesWhenMissing(t *testing.T) {
	fake := cloud.NewFake()

	n, err := EnsureNetwork(context.Background(), fake, "shoot-net", "10.0.0.0/16")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil || n.Name != "shoot-net" || n.CIDR != "10.0.0.0/16" {
		t.Fatalf("got network %+v", n)
	}
	if want := []string{"get:shoot-net", "create:shoot-net"}; !slices.Equal(fake.Calls(), want) {
		t.Errorf("calls = %v, want %v", fake.Calls(), want)
	}
}

func TestEnsureNetwork_IsIdempotent(t *testing.T) {
	fake := cloud.NewFake()
	ctx := context.Background()

	for range 3 {
		if _, err := EnsureNetwork(ctx, fake, "shoot-net", "10.0.0.0/16"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	creates := 0
	for _, c := range fake.Calls() {
		if c == "create:shoot-net" {
			creates++
		}
	}
	if creates != 1 {
		t.Errorf("network was created %d times, want exactly 1 (calls: %v)", creates, fake.Calls())
	}
}

func TestEnsureNetwork_CIDRMismatch(t *testing.T) {
	fake := cloud.NewFake(cloud.Network{Name: "shoot-net", CIDR: "192.168.0.0/24"})

	_, err := EnsureNetwork(context.Background(), fake, "shoot-net", "10.0.0.0/16")
	if !errors.Is(err, ErrCIDRMismatch) {
		t.Fatalf("err = %v, want it to wrap ErrCIDRMismatch", err)
	}
	if slices.Contains(fake.Calls(), "create:shoot-net") {
		t.Errorf("must not try to create an existing network")
	}
}

func TestEnsureNetwork_WrapsCloudErrors(t *testing.T) {
	t.Run("get fails", func(t *testing.T) {
		fake := cloud.NewFake()
		fake.GetErr = errBoom
		_, err := EnsureNetwork(context.Background(), fake, "shoot-net", "10.0.0.0/16")
		if !errors.Is(err, errBoom) {
			t.Fatalf("err = %v, want it to wrap the cloud error", err)
		}
	})
	t.Run("create fails", func(t *testing.T) {
		fake := cloud.NewFake()
		fake.CreateErr = errBoom
		_, err := EnsureNetwork(context.Background(), fake, "shoot-net", "10.0.0.0/16")
		if !errors.Is(err, errBoom) {
			t.Fatalf("err = %v, want it to wrap the cloud error", err)
		}
	})
}

func TestDeleteNetwork(t *testing.T) {
	ctx := context.Background()

	t.Run("deletes existing network", func(t *testing.T) {
		fake := cloud.NewFake(cloud.Network{Name: "shoot-net", CIDR: "10.0.0.0/16"})
		if err := DeleteNetwork(ctx, fake, "shoot-net"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := fake.Networks()["shoot-net"]; ok {
			t.Errorf("network still exists")
		}
	})
	t.Run("already gone is fine", func(t *testing.T) {
		if err := DeleteNetwork(ctx, cloud.NewFake(), "shoot-net"); err != nil {
			t.Fatalf("deleting a missing network must not fail, got: %v", err)
		}
	})
	t.Run("real errors are returned", func(t *testing.T) {
		fake := cloud.NewFake(cloud.Network{Name: "shoot-net"})
		fake.DeleteErr = errBoom
		if err := DeleteNetwork(ctx, fake, "shoot-net"); !errors.Is(err, errBoom) {
			t.Fatalf("err = %v, want it to wrap the cloud error", err)
		}
	})
}
