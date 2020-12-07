// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by deepcopy-gen. DO NOT EDIT.

package responses

import (
	models "github.com/SolarLabRU/fastpay-go-commons/models"
	bus "github.com/SolarLabRU/fastpay-go-commons/remotes/bus"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccountBalanceData) DeepCopyInto(out *AccountBalanceData) {
	*out = *in
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]models.AmountOfBank, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccountBalanceData.
func (in *AccountBalanceData) DeepCopy() *AccountBalanceData {
	if in == nil {
		return nil
	}
	out := new(AccountBalanceData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccountBalanceResponse) DeepCopyInto(out *AccountBalanceResponse) {
	*out = *in
	in.Data.DeepCopyInto(&out.Data)
	out.BaseResponse = in.BaseResponse
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccountBalanceResponse.
func (in *AccountBalanceResponse) DeepCopy() *AccountBalanceResponse {
	if in == nil {
		return nil
	}
	out := new(AccountBalanceResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyResponse is an autogenerated deepcopy function, copying the receiver, creating a new bus.Response.
func (in *AccountBalanceResponse) DeepCopyResponse() bus.Response {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccountResponse) DeepCopyInto(out *AccountResponse) {
	*out = *in
	in.Data.DeepCopyInto(&out.Data)
	out.BaseResponse = in.BaseResponse
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccountResponse.
func (in *AccountResponse) DeepCopy() *AccountResponse {
	if in == nil {
		return nil
	}
	out := new(AccountResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyResponse is an autogenerated deepcopy function, copying the receiver, creating a new bus.Response.
func (in *AccountResponse) DeepCopyResponse() bus.Response {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
