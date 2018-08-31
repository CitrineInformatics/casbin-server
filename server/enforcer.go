// Copyright 2018 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"github.com/casbin/casbin"
	pb "github.com/casbin/casbin-server/proto"
	"github.com/casbin/casbin/persist"
	"github.com/casbin/casbin/persist/file-adapter"
	"golang.org/x/net/context"
)

// Server is used to implement proto.CasbinServer.
type Server struct{
	enforcerMap map[int]casbin.Enforcer
	adapterMap  map[int]persist.Adapter
}

func (s *Server) getAdapter(handle int) persist.Adapter {
	if _, ok := s.adapterMap[handle]; ok {
		return s.adapterMap[handle]
	} else {
		return nil
	}
}

func (s *Server) addEnforcer(e casbin.Enforcer) int {
	cnt := len(s.enforcerMap)
	s.enforcerMap[cnt] = e
	return cnt
}

func (s *Server) addAdapter(a persist.Adapter) int {
	cnt := len(s.adapterMap)
	s.adapterMap[cnt] = a
	return cnt
}

func (s *Server) NewEnforcer(ctx context.Context, in *pb.NewEnforcerRequest) (*pb.NewEnforcerReply, error) {
	a := s.getAdapter(int(in.AdapterHandle))
	e := casbin.NewEnforcer(in.ModelText, a)

	h := s.addEnforcer(*e)

	return &pb.NewEnforcerReply{Handler: int32(h)}, nil
}

func (s *Server) NewAdapter(ctx context.Context, in *pb.NewAdapterRequest) (*pb.NewAdapterReply, error) {
	a := fileadapter.NewAdapter(in.ConnectString)

	h := s.addAdapter(a)

	return &pb.NewAdapterReply{Handler: int32(h)}, nil
}