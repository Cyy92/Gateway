# Generated by the protocol buffer compiler.  DO NOT EDIT!
# Source: fxwatcher.proto for package 'pb'

require 'grpc'
require 'fxwatcher_pb'

module Pb
  module FxWatcher
    class Service

      include GRPC::GenericService

      self.marshal_class_method = :encode
      self.unmarshal_class_method = :decode
      self.service_name = 'pb.FxWatcher'

      rpc :Call, Request, Reply
    end

    Stub = Service.rpc_stub_class
  end
end
