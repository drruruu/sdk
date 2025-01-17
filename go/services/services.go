package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/trinsic-id/okapi/go/okapi"
	"github.com/trinsic-id/okapi/go/okapiproto"
	sdk "github.com/trinsic-id/sdk/go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type Document map[string]interface{}

type ServiceBase struct {
	capabilityInvocation string
}

type Service interface {
	GetMetadataContext(userContext context.Context) (context.Context, error)
	GetMetadata() (metadata.MD, error)
	SetProfile(profile *sdk.WalletProfile) error
}

func (s *ServiceBase) GetMetadataContext(userContext context.Context) (context.Context, error) {
	md, err := s.GetMetadata()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		return nil, errors.New("userContext cannot be nil")
	}
	return metadata.NewOutgoingContext(userContext, md), nil
}

func (s *ServiceBase) GetMetadata() (metadata.MD, error) {
	if s.capabilityInvocation == "" {
		return nil, errors.New("profile not set")
	}
	return metadata.New(map[string]string{
		"capability-invocation": s.capabilityInvocation,
	}), nil
}

func (s *ServiceBase) SetProfile(profile *sdk.WalletProfile) error {
	capabilityStruct, err := structpb.NewStruct(map[string]interface{}{
		"@context":         "https://w3id.org/security/v2",
		"invocationTarget": profile.WalletId,
		"proof": map[string]interface{}{
			"proofPurpose": "capabilityInvocation",
			"created":      time.Now().Format(time.RFC3339),
			"capability":   profile.Capability,
		},
	})
	if err != nil {
		return err
	}

	invokerKey := okapiproto.JsonWebKey{}
	err = proto.Unmarshal(profile.InvokerJwk, &invokerKey)
	if err != nil {
		return err
	}

	proofResponse, err := okapi.LdProofs().CreateProof(&okapiproto.CreateProofRequest{
		Document: capabilityStruct,
		Key:      &invokerKey,
		Suite:    okapiproto.LdSuite_LD_SUITE_JCSED25519SIGNATURE2020,
	})
	if err != nil {
		return err
	}

	proofJson, err := json.Marshal(proofResponse.SignedDocument.AsMap())
	if err != nil {
		return err
	}

	s.capabilityInvocation = base64.StdEncoding.EncodeToString(proofJson)
	return nil
}

type WalletService interface {
	Service
	RegisterOrConnect(userContext context.Context, email string) error
	CreateWallet(userContext context.Context, securityCode string) (*sdk.WalletProfile, error)
	IssueCredential(userContext context.Context, document Document) (Document, error)
	Search(userContext context.Context, query string) (*sdk.SearchResponse, error)
	InsertItem(userContext context.Context, item Document) (string, error)
	Send(userContext context.Context, document Document, email string) error
	CreateProof(userContext context.Context, documentId string, revealDocument Document) (Document, error)
	VerifyProof(userContext context.Context, proofDocument Document) (bool, error)
}

func CreateWalletService(serviceAddress string, channel *grpc.ClientConn) (WalletService, error) {
	channel, err := CreateChannelIfNeeded(serviceAddress, channel, true)
	if err != nil {
		return nil, err
	}

	service := &WalletBase{
		ServiceBase:      &ServiceBase{},
		channel:          channel,
		walletClient:     sdk.NewWalletClient(channel),
		credentialClient: sdk.NewCredentialClient(channel),
	}

	return service, nil
}

func CreateChannelIfNeeded(serviceAddress string, channel *grpc.ClientConn, blockOnOpen bool) (*grpc.ClientConn, error) {
	if channel == nil {
		var serviceUrl, err = url.Parse(serviceAddress)
		if err != nil {
			return nil, err
		}
		if serviceUrl.Port() == "" {
			return nil, &url.Error{Op: "parse", URL: serviceAddress, Err: errors.New("missing port (or scheme) in URL")}
		}
		dialUrl := serviceUrl.Hostname() + ":" + serviceUrl.Port()
		var dialOptions []grpc.DialOption
		if blockOnOpen {
			dialOptions = append(dialOptions, grpc.WithBlock())
		}
		if serviceUrl.Scheme == "http" {
			dialOptions = append(dialOptions, grpc.WithInsecure())
		} else {
			// TODO - Get the credentials bundle
			//credBundle := credentials.Bundle{}
			//dialOptions = append(dialOptions, grpc.WithCredentialsBundle(credentials.Bundle()))
			return nil, errors.New("HTTPS not supported yet due to credential bundle declaration")
		}
		channel, err = grpc.Dial(dialUrl, dialOptions...)
		if err != nil {
			return nil, err
		}
	}
	return channel, nil
}

type WalletBase struct {
	*ServiceBase
	channel          *grpc.ClientConn
	walletClient     sdk.WalletClient
	credentialClient sdk.CredentialClient
}

func (w *WalletBase) RegisterOrConnect(userContext context.Context, email string) error {
	connectRequest := sdk.ConnectRequest{
		ContactMethod: &sdk.ConnectRequest_Email{Email: email},
	}

	md, err := w.GetMetadataContext(userContext)
	if err != nil {
		return err
	}
	_, err = w.walletClient.ConnectExternalIdentity(md, &connectRequest)
	if err != nil {
		return err
	}
	return nil
}

func (w *WalletBase) CreateWallet(userContext context.Context, securityCode string) (*sdk.WalletProfile, error) {
	dk := okapi.DidKey()

	// Generate new DID used by the current device
	myKey, err := dk.Generate(&okapiproto.GenerateKeyRequest{KeyType: okapiproto.KeyType_KEY_TYPE_ED25519})
	if err != nil {
		return nil, err
	}

	myDidDocument := myKey.DidDocument.AsMap()
	walletRequest := &sdk.CreateWalletRequest{
		Controller:   myDidDocument["id"].(string),
		SecurityCode: securityCode,
	}

	createWalletResponse, err := w.walletClient.CreateWallet(userContext, walletRequest)
	if err != nil {
		return nil, err
	}

	keyBytes, err := proto.Marshal(myKey.Key[0])
	if err != nil {
		return nil, err
	}
	return &sdk.WalletProfile{
		DidDocument: &sdk.JsonPayload{Json: &sdk.JsonPayload_JsonStruct{JsonStruct: myKey.DidDocument}},
		WalletId:    createWalletResponse.WalletId,
		Invoker:     createWalletResponse.Invoker,
		Capability:  createWalletResponse.Capability,
		InvokerJwk:  keyBytes,
	}, nil
}

func (w *WalletBase) IssueCredential(userContext context.Context, document Document) (Document, error) {
	jsonBytes, err := json.Marshal(document)
	if err != nil {
		return nil, err
	}
	issueRequest := sdk.IssueRequest{
		Document: &sdk.JsonPayload{
			Json: &sdk.JsonPayload_JsonString{
				JsonString: string(jsonBytes),
			},
		},
	}

	md, err := w.GetMetadataContext(userContext)
	if err != nil {
		return nil, err
	}
	response, err := w.credentialClient.Issue(md, &issueRequest)
	if err != nil {
		return nil, err
	}
	var doc map[string]interface{}
	err = json.Unmarshal([]byte(response.Document.GetJsonString()), &doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (w *WalletBase) Search(userContext context.Context, query string) (*sdk.SearchResponse, error) {
	md, err := w.GetMetadataContext(userContext)
	if err != nil {
		return nil, err
	}
	response, err := w.walletClient.Search(md, &sdk.SearchRequest{
		Query: query,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (w *WalletBase) InsertItem(userContext context.Context, item Document) (string, error) {
	jsonString, err := json.Marshal(item)
	if err != nil {
		return "", err
	}
	md, err := w.GetMetadataContext(userContext)
	if err != nil {
		return "", err
	}
	response, err := w.walletClient.InsertItem(md, &sdk.InsertItemRequest{
		Item: &sdk.JsonPayload{
			Json: &sdk.JsonPayload_JsonString{
				JsonString: string(jsonString),
			},
		},
	})
	if err != nil {
		return "", err
	}
	return response.ItemId, nil
}

func (w *WalletBase) Send(userContext context.Context, document Document, email string) error {
	jsonString, err := json.Marshal(document)
	if err != nil {
		return err
	}
	md, err := w.GetMetadataContext(userContext)
	if err != nil {
		return err
	}
	_, err = w.credentialClient.Send(md, &sdk.SendRequest{
		DeliveryMethod: &sdk.SendRequest_Email{
			Email: email,
		},
		Document: &sdk.JsonPayload{
			Json: &sdk.JsonPayload_JsonString{
				JsonString: string(jsonString),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (w *WalletBase) CreateProof(userContext context.Context, documentId string, revealDocument Document) (Document, error) {
	jsonString, err := json.Marshal(revealDocument)
	if err != nil {
		return nil, err
	}
	md, err := w.GetMetadataContext(userContext)
	if err != nil {
		return nil, err
	}
	proof, err := w.credentialClient.CreateProof(md, &sdk.CreateProofRequest{
		DocumentId: documentId,
		RevealDocument: &sdk.JsonPayload{
			Json: &sdk.JsonPayload_JsonString{
				JsonString: string(jsonString),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	var proofMap map[string]interface{}
	err = json.Unmarshal([]byte(proof.ProofDocument.GetJsonString()), &proofMap)
	if err != nil {
		return nil, err
	}
	return proofMap, nil
}

func (w *WalletBase) VerifyProof(userContext context.Context, proofDocument Document) (bool, error) {
	jsonString, err := json.Marshal(proofDocument)
	if err != nil {
		return false, err
	}
	md, err := w.GetMetadataContext(userContext)
	if err != nil {
		return false, err
	}
	proof, err := w.credentialClient.VerifyProof(md, &sdk.VerifyProofRequest{
		ProofDocument: &sdk.JsonPayload{
			Json: &sdk.JsonPayload_JsonString{
				JsonString: string(jsonString),
			},
		},
	})
	if err != nil {
		return false, err
	}
	return proof.Valid, nil
}

type ProviderService interface {
	Service
	InviteParticipant(userContext context.Context, request *sdk.InviteRequest) (*sdk.InviteResponse, error)
	InvitationStatus(userContext context.Context, request *sdk.InvitationStatusRequest) (*sdk.InvitationStatusResponse, error)
}

type ProviderBase struct {
	*ServiceBase
	channel        *grpc.ClientConn
	providerClient sdk.ProviderClient
}

func CreateProviderService(serviceAddress string, channel *grpc.ClientConn) (ProviderService, error) {
	channel, err := CreateChannelIfNeeded(serviceAddress, channel, true)
	if err != nil {
		return nil, err
	}

	service := &ProviderBase{
		ServiceBase:    &ServiceBase{},
		channel:        channel,
		providerClient: sdk.NewProviderClient(channel),
	}
	return service, nil
}

func (p *ProviderBase) InviteParticipant(userContext context.Context, request *sdk.InviteRequest) (*sdk.InviteResponse, error) {
	// Verify contact method is set
	switch request.ContactMethod.(type) {
	case nil:
		return nil, fmt.Errorf("unset contact method")
	}
	response, err := p.providerClient.Invite(userContext, request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (p *ProviderBase) InvitationStatus(userContext context.Context, request *sdk.InvitationStatusRequest) (*sdk.InvitationStatusResponse, error) {
	response, err := p.providerClient.InvitationStatus(userContext, request)
	if err != nil {
		return nil, err
	}
	return response, nil
}
