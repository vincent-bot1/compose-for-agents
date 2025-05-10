package main

import "context"

type Secret struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func secretValue(ctx context.Context, id string) (string, error) {
	var secret Secret
	if err := get(ctx, httpClient(dialJFS), "/secrets/"+id, &secret); err != nil {
		return "", err
	}

	return secret.Value, nil
}
