#[derive(Clone, PartialEq, ::prost::Message)]
pub struct InviteRequest {
    #[prost(enumeration = "ParticipantType", tag = "1")]
    pub participant: i32,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(oneof = "invite_request::ContactMethod", tags = "5, 6, 7")]
    pub contact_method: ::core::option::Option<invite_request::ContactMethod>,
}
/// Nested message and enum types in `InviteRequest`.
pub mod invite_request {
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct DidCommInvitation {}
    #[derive(Clone, PartialEq, ::prost::Oneof)]
    pub enum ContactMethod {
        #[prost(string, tag = "5")]
        Email(::prost::alloc::string::String),
        #[prost(string, tag = "6")]
        Phone(::prost::alloc::string::String),
        #[prost(message, tag = "7")]
        DidcommInvitation(DidCommInvitation),
    }
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct InviteResponse {
    #[prost(enumeration = "super::super::common::v1::ResponseStatus", tag = "1")]
    pub status: i32,
    #[prost(string, tag = "10")]
    pub invitation_id: ::prost::alloc::string::String,
}
/// Request details for the status of onboarding
/// an individual or organization.
/// The referenece_id passed is the response from the
/// `Onboard` method call
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct InvitationStatusRequest {
    #[prost(string, tag = "1")]
    pub invitation_id: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct InvitationStatusResponse {
    #[prost(enumeration = "invitation_status_response::Status", tag = "1")]
    pub status: i32,
    #[prost(string, tag = "2")]
    pub status_details: ::prost::alloc::string::String,
}
/// Nested message and enum types in `InvitationStatusResponse`.
pub mod invitation_status_response {
    #[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
    #[repr(i32)]
    pub enum Status {
        /// Onboarding resulted in error
        Error = 0,
        /// The participant has been invited
        InvitationSent = 1,
        /// The participant has been onboarded
        Completed = 2,
    }
}
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum ParticipantType {
    Individual = 0,
    Organization = 1,
}
#[doc = r" Generated client implementations."]
pub mod provider_client {
    #![allow(unused_variables, dead_code, missing_docs)]
    use tonic::codegen::*;
    pub struct ProviderClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl ProviderClient<tonic::transport::Channel> {
        #[doc = r" Attempt to create a new client by connecting to a given endpoint."]
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: std::convert::TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> ProviderClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::BoxBody>,
        T::ResponseBody: Body + HttpBody + Send + 'static,
        T::Error: Into<StdError>,
        <T::ResponseBody as HttpBody>::Error: Into<StdError> + Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_interceptor(inner: T, interceptor: impl Into<tonic::Interceptor>) -> Self {
            let inner = tonic::client::Grpc::with_interceptor(inner, interceptor);
            Self { inner }
        }
        #[doc = "   rpc CreateOrganization(CreateOrganizationRequest) returns (CreateOrganizationResponse);"]
        pub async fn invite(
            &mut self,
            request: impl tonic::IntoRequest<super::InviteRequest>,
        ) -> Result<tonic::Response<super::InviteResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/services.provider.v1.Provider/Invite");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn invite_with_workflow(
            &mut self,
            request: impl tonic::IntoRequest<super::InviteRequest>,
        ) -> Result<tonic::Response<super::InviteResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/services.provider.v1.Provider/InviteWithWorkflow",
            );
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn invitation_status(
            &mut self,
            request: impl tonic::IntoRequest<super::InvitationStatusRequest>,
        ) -> Result<tonic::Response<super::InvitationStatusResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/services.provider.v1.Provider/InvitationStatus",
            );
            self.inner.unary(request.into_request(), path, codec).await
        }
    }
    impl<T: Clone> Clone for ProviderClient<T> {
        fn clone(&self) -> Self {
            Self {
                inner: self.inner.clone(),
            }
        }
    }
    impl<T> std::fmt::Debug for ProviderClient<T> {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "ProviderClient {{ ... }}")
        }
    }
}