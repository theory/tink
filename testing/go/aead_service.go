// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
///////////////////////////////////////////////////////////////////////////////

package services

import (
	"context"
	"fmt"

	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/aead/internal/testing/kmsaead"
	"github.com/google/tink/go/core/registry"
	pb "github.com/google/tink/testing/go/protos/testing_api_go_grpc"
)

// AEADService implements the Aead testing service.
type AEADService struct {
	pb.AeadServer
}

func (s *AEADService) Create(ctx context.Context, req *pb.CreationRequest) (*pb.CreationResponse, error) {
	handle, err := toKeysetHandle(req.GetAnnotatedKeyset())
	if err != nil {
		return &pb.CreationResponse{Err: err.Error()}, nil
	}
	_, err = aead.New(handle)
	if err != nil {
		return &pb.CreationResponse{Err: err.Error()}, nil
	}
	return &pb.CreationResponse{}, nil
}

func (s *AEADService) Encrypt(ctx context.Context, req *pb.AeadEncryptRequest) (*pb.AeadEncryptResponse, error) {
	handle, err := toKeysetHandle(req.GetAnnotatedKeyset())
	if err != nil {
		return nil, err
	}
	cipher, err := aead.New(handle)
	if err != nil {
		return nil, err
	}
	ciphertext, err := cipher.Encrypt(req.Plaintext, req.AssociatedData)
	if err != nil {
		return &pb.AeadEncryptResponse{
			Result: &pb.AeadEncryptResponse_Err{err.Error()}}, nil
	}
	return &pb.AeadEncryptResponse{
		Result: &pb.AeadEncryptResponse_Ciphertext{ciphertext}}, nil
}

func (s *AEADService) Decrypt(ctx context.Context, req *pb.AeadDecryptRequest) (*pb.AeadDecryptResponse, error) {
	handle, err := toKeysetHandle(req.GetAnnotatedKeyset())
	if err != nil {
		return nil, err
	}
	cipher, err := aead.New(handle)
	if err != nil {
		return nil, err
	}
	plaintext, err := cipher.Decrypt(req.Ciphertext, req.AssociatedData)
	if err != nil {
		return &pb.AeadDecryptResponse{
			Result: &pb.AeadDecryptResponse_Err{err.Error()}}, nil
	}
	return &pb.AeadDecryptResponse{
		Result: &pb.AeadDecryptResponse_Plaintext{plaintext}}, nil
}

func init() {
	if err := registry.RegisterKeyManager(kmsaead.NewKeyManager()); err != nil {
		panic(fmt.Sprintf("registry.RegisterKeyManager(kmsaead.NewKeyManager()) failed: %v", err))
	}
}
