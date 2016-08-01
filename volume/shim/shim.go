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

func (d *shimDriver) Create(req volumeplugin.Request) volumeplugin.Response {
	var res volumeplugin.Response
	_, err := d.d.Create(req.Name, req.Options)
	if err != nil {
		res.Err = err.Error()
	}
	return res
}

func (d *shimDriver) List(req volumeplugin.Request) volumeplugin.Response {
	var res volumeplugin.Response
	ls, err := d.d.List()
	if err != nil {
		res.Err = err.Error()
		return res
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
	return res
}

func (d *shimDriver) Get(req volumeplugin.Request) volumeplugin.Response {
	var res volumeplugin.Response
	v, err := d.d.Get(req.Name)
	if err != nil {
		res.Err = err.Error()
		return res
	}
	res.Volume = &volumeplugin.Volume{
		Name:       v.Name(),
		Mountpoint: v.Path(),
		Status:     v.Status(),
	}
	return res
}

func (d *shimDriver) Remove(req volumeplugin.Request) volumeplugin.Response {
	var res volumeplugin.Response
	v, err := d.d.Get(req.Name)
	if err != nil {
		res.Err = err.Error()
		return res
	}
	if err := d.d.Remove(v); err != nil {
		res.Err = err.Error()
	}
	return res
}

func (d *shimDriver) Path(req volumeplugin.Request) volumeplugin.Response {
	var res volumeplugin.Response
	v, err := d.d.Get(req.Name)
	if err != nil {
		res.Err = err.Error()
		return res
	}
	res.Mountpoint = v.Path()
	return res
}

func (d *shimDriver) Mount(req volumeplugin.MountRequest) volumeplugin.Response {
	var res volumeplugin.Response
	v, err := d.d.Get(req.Name)
	if err != nil {
		res.Err = err.Error()
		return res
	}
	pth, err := v.Mount(req.ID)
	if err != nil {
		res.Err = err.Error()
	}
	res.Mountpoint = pth
	return res
}

func (d *shimDriver) Unmount(req volumeplugin.UnmountRequest) volumeplugin.Response {
	var res volumeplugin.Response
	v, err := d.d.Get(req.Name)
	if err != nil {
		res.Err = err.Error()
		return res
	}
	if err := v.Unmount(req.ID); err != nil {
		res.Err = err.Error()
	}
	return res
}

func (d *shimDriver) Capabilities(req volumeplugin.Request) volumeplugin.Response {
	var res volumeplugin.Response
	res.Capabilities = volumeplugin.Capability{Scope: d.d.Scope()}
	return res
}
