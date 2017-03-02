package shim

import (
	"github.com/docker/docker/volume"
	volumeplugin "github.com/docker/go-plugins-helpers/volume"
)

type shimDriver struct {
	d volume.Driver
}

// NewHandlerFromVolumeDriver creates a plugin handler from an existing volume
// driver. This could be used, for instance, by the `local` volume driver built-in
// to Docker Engine and it would create a plugin from it that maps plugin API calls
// directly to any volume driver that satifies the volume.Driver interface from
// Docker Engine.
func NewHandlerFromVolumeDriver(d volume.Driver) *volumeplugin.Handler {
	return volumeplugin.NewHandler(&shimDriver{d})
}

func (d *shimDriver) Create(req *volumeplugin.CreateRequest) error {
	_, err := d.d.Create(req.Name, req.Options)
	return err
}

func (d *shimDriver) List() (*volumeplugin.ListResponse, error) {
	var res *volumeplugin.ListResponse
	ls, err := d.d.List()
	if err != nil {
		res.Err = err.Error()
		return &volumeplugin.ListResponse{}, err
	}
	vols := make([]*volumeplugin.Volume, len(ls))

	for i, v := range ls {
		vol := &volumeplugin.Volume{
			Name:       v.Name(),
			Mountpoint: v.Path(),
		}
		vols[i] = vol
	}
	res.Volumes = vols
	return res, nil
}

func (d *shimDriver) Get(req *volumeplugin.GetRequest) (*volumeplugin.GetResponse, error) {
	var res *volumeplugin.GetResponse
	v, err := d.d.Get(req.Name)
	if err != nil {
		res.Err = err.Error()
		return &volumeplugin.GetResponse{}, err
	}
	res.Volume = &volumeplugin.Volume{
		Name:       v.Name(),
		Mountpoint: v.Path(),
		Status:     v.Status(),
	}
	return &volumeplugin.GetResponse{}, nil
}

func (d *shimDriver) Remove(req *volumeplugin.RemoveRequest) error {
	v, err := d.d.Get(req.Name)
	if err != nil {
		return err
	}
	if err := d.d.Remove(v); err != nil {
		return err
	}
	return nil
}

func (d *shimDriver) Path(req *volumeplugin.PathRequest) (*volumeplugin.PathResponse, error) {
	var res *volumeplugin.PathResponse
	v, err := d.d.Get(req.Name)
	if err != nil {
		res.Err = err.Error()
		return res, err
	}
	res.Mountpoint = v.Path()
	return res, nil
}

func (d *shimDriver) Mount(req *volumeplugin.MountRequest) (*volumeplugin.MountResponse, error) {
	var res *volumeplugin.MountResponse
	v, err := d.d.Get(req.Name)
	if err != nil {
		res.Err = err.Error()
		return res, err
	}
	pth, err := v.Mount(req.ID)
	if err != nil {
		res.Err = err.Error()
	}
	res.Mountpoint = pth
	return res, nil
}

func (d *shimDriver) Unmount(req *volumeplugin.UnmountRequest) error {
	v, err := d.d.Get(req.Name)
	if err != nil {
		return err
	}
	if err := v.Unmount(req.ID); err != nil {
		return err
	}
	return nil
}

func (d *shimDriver) Capabilities() *volumeplugin.CapabilitiesResponse {
	var res *volumeplugin.CapabilitiesResponse
	res.Capabilities = volumeplugin.Capability{Scope: d.d.Scope()}
	return res
}
