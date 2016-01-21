package ltm

import (
	"github.com/zalando-techmonkeys/baboon-proxy/backend"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
	"github.com/zalando-techmonkeys/baboon-proxy/errors"
	"net/url"
	"path"
)

/*
{"kind":"tm:ltm:datagroup:datagroupcollectionstate","selfLink":"https://localhost/mgmt/tm/ltm/datagroup?ver=11.5.1","items":[{}]}
*/

// DataGroups struct provides information
// about multiple datagroups
type DataGroups struct {
	Kind  string `json:"kind"`
	Items []struct {
		DataGroup
	} `json:"items"`
}

// DataGroup struct provides information
// about a specific datagroup
type DataGroup struct {
	Kind       string `json:"kind"`
	Name       string `json:"name"`
	Partition  string `json:"partition"`
	FullPath   string `json:"fullPath"`
	Generation int    `json:"generation"`
	Type       string `json:"type"`
	Records    []struct {
		Name string `json:"name"`
		Data string `json:"data"`
	} `json:"records"`
}

// CreateDataGroup struct to create a datagroup
type CreateDataGroup struct {
	Name    string `json:"name" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Records []struct {
		Name string `json:"name" binding:"required"`
		Data string `json:"data" binding:"required"`
	} `json:"records" binding:"required"`
}

// CreateDataGroupItem struct to add records in a existing datagroup
type CreateDataGroupItem struct {
	Records []struct {
		Name string `json:"name"`
		Data string `json:"data"`
	} `json:"records" binding:"required"`
}

// RemoveDataGroup for deleting datagroup
type RemoveDataGroup struct {
	Name string `json:"name" binding:"required"`
}

// ShowLTMDataGroup lists all datagroups on a loadbalancer
func ShowLTMDataGroup(host, source string) (*backend.Response, *DataGroups, *errors.Error) {
	// Declaration LTM DataGroup
	ltmdatagroup := new(DataGroups)
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, common.Dg, source)
	res, err := backend.Request(common.GET, u.String(), ltmdatagroup)
	if err != nil {
		return nil, nil, err
	}
	return res, ltmdatagroup, nil
}

// ShowLTMDataGroupName lists a specific datagroup on a loadbalancer
func ShowLTMDataGroupName(host, direction, datagroupname string) (*backend.Response, *DataGroup, *errors.Error) {
	// Declaration LTM DataGroup by Name
	ltmdatagroupname := new(DataGroup)
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, common.Dg, direction, "/", datagroupname)
	res, err := backend.Request(common.GET, u.String(), ltmdatagroupname)
	if err != nil {
		return nil, nil, err
	}
	return res, ltmdatagroupname, nil
}

// PostLTMDataGroup creates a new datagroup
func PostLTMDataGroup(host, direction string, json *CreateDataGroup) (*backend.Response, *errors.Error) {
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, common.Dg, direction)
	r, err := backend.Request(common.POST, u.String(), &json)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// PutLTMDataGroupItem deletes all records in a datagroup and add new records
func PutLTMDataGroupItem(host, direction, datagroup string, json *CreateDataGroupItem) (*backend.Response, *errors.Error) {
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, common.Dg, direction, "/", datagroup)

	r, err := backend.Request(common.PATCH, u.String(), &json)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// PatchLTMDataGroupItem keeps all records in a datagroup and add new records
// F5 API use PATCH and PUT in the same way. Overwriting all records, which is bad if you want to add items
// in existing list. It gets all records first an append the new records from client.
func PatchLTMDataGroupItem(host, direction, datagroup string, json *CreateDataGroupItem) (*backend.Response, *errors.Error) {
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, common.Dg, direction, "/", datagroup)

	_, data, err := ShowLTMDataGroupName(host, direction, datagroup)
	if err != nil {
		return nil, err
	}
	for _, v := range data.Records {
		json.Records = append(json.Records, v)
	}
	r, err := backend.Request(common.PATCH, u.String(), &json)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// DeleteLTMDataGroup deletes a specific datagroup
func DeleteLTMDataGroup(host, direction, datagroupname string) (*backend.Response, *errors.Error) {
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, common.Dg, direction, "/", datagroupname)
	r, err := backend.Request(common.DELETE, u.String(), nil)
	if err != nil {
		return nil, err
	}
	return r, nil
}
