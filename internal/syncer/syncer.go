package syncer

import (
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// Copy 镜像同步
func Copy(
	src string,
	dst string,
	srcKeychain authn.Keychain,
	dstKeychain authn.Keychain,
) error {

	return CopyWithPlatform(
		src,
		dst,
		"",
		srcKeychain,
		dstKeychain,
	)
}

// CopyWithPlatform
//
// source 使用 srcKeychain
// destination 使用 dstKeychain
func CopyWithPlatform(
	src string,
	dst string,
	targetPlatform string,
	srcKeychain authn.Keychain,
	dstKeychain authn.Keychain,
) error {

	srcRef, err := name.ParseReference(
		src,
	)

	if err != nil {
		return err
	}

	dstRef, err := name.ParseReference(
		dst,
	)

	if err != nil {
		return err
	}

	srcOptions := []remote.Option{}

	if srcKeychain != nil {

		srcOptions = append(
			srcOptions,
			remote.WithAuthFromKeychain(
				srcKeychain,
			),
		)
	}

	if targetPlatform != "" {

		srcOptions = append(
			srcOptions,
			remote.WithPlatform(
				parsePlatform(targetPlatform),
			),
		)
	}

	img, err := remote.Image(
		srcRef,
		srcOptions...,
	)

	if err != nil {
		return err
	}

	dstOptions := []remote.Option{}

	if dstKeychain != nil {

		dstOptions = append(
			dstOptions,
			remote.WithAuthFromKeychain(
				dstKeychain,
			),
		)
	}

	return remote.Write(
		dstRef,
		img,
		dstOptions...,
	)
}

// parsePlatform
//
// 格式:
// linux/amd64
// linux/arm64/v8
func parsePlatform(
	value string,
) v1.Platform {

	parts := strings.Split(
		value,
		"/",
	)

	p := v1.Platform{}

	if len(parts) > 0 {
		p.OS = parts[0]
	}

	if len(parts) > 1 {
		p.Architecture = parts[1]
	}

	if len(parts) > 2 {
		p.Variant = parts[2]
	}

	return p
}
