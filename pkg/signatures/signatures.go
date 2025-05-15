package signatures

import (
	"context"
	"crypto"
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/sigstore/cosign/v2/pkg/cosign"
	ociremote "github.com/sigstore/cosign/v2/pkg/oci/remote"
	"github.com/sigstore/cosign/v2/pkg/signature"
	"golang.org/x/sync/errgroup"
)

const keyRef = "https://raw.githubusercontent.com/docker/keyring/refs/heads/main/public/mcp/latest.pub"

func Verify(ctx context.Context, images []string) error {
	signatures, err := name.NewRepository("mcp/signatures")
	if err != nil {
		return err
	}

	rekor, err := cosign.GetRekorPubs(ctx)
	if err != nil {
		return fmt.Errorf("getting Rekor public keys: %w", err)
	}

	key, err := signature.PublicKeyFromKeyRefWithHashAlgo(ctx, keyRef, crypto.SHA256)
	if err != nil {
		return fmt.Errorf("loading public key: %w", err)
	}

	errs, ctxVerify := errgroup.WithContext(ctx)
	errs.SetLimit(2)
	for _, img := range images {
		img := img
		errs.Go(func() error {
			ref, err := name.ParseReference(img)
			if err != nil {
				return fmt.Errorf("parsing reference: %w", err)
			}

			_, bundleVerified, err := cosign.VerifyImageSignatures(ctx, ref, &cosign.CheckOpts{
				RegistryClientOpts: []ociremote.Option{
					ociremote.WithTargetRepository(signatures),
					ociremote.WithRemoteOptions(
						remote.WithContext(ctxVerify),
						remote.WithUserAgent("docker/mcp_gateway"),
						// remote.WithAuthFromKeychain(authn.DefaultKeychain),
					),
				},
				RekorPubKeys: rekor,
				SigVerifier:  key,
			})
			if err != nil {
				return err
			}

			if !bundleVerified {
				return fmt.Errorf("bundle verification failed")
			}

			return nil
		})
	}

	return errs.Wait()
}
