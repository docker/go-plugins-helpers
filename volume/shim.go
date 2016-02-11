package volume

import "github.com/docker/docker/volume"

type shimDriver struct {
	d volume.Driver
}

// NewHandlerFromVolumeDriver creates a plugin handler from an existing volume
// driver. This could be used, for instance, by the `local` volume driver built-in
// to Docker Engine and it would create a plugin from it that maps plugin API calls
// directly to any volume driver that satifies the volume.Driver interface from
// Docker Engine.
func NewHandlerFromVolumeDriver(d volume.Driver) *Handler {
	return NewHandler(&shimDriver{d})
}

func (d *shimDriver) Create(req Request) Response {
	var res Response
	_, err := d.d.Create(req.Name, req.Options)
	if err != nil {
		res.Err = err.Error()
	}
	return res
}

func (d *shimDriver) List(req Request) Response {
	var res Response
	ls, err := d.d.List()
	if err != nil {
		res.Err = err.Error()
		return res
	}
	vols := make([]*Volume, len(ls))

	for _, v := range ls {
		vol := &Volume{
			Name:       v.Name(),
			Mountpoint: v.Path(),
		}
		vols = append(vols, vol)
	}
	res.Volumes = vols
	return res
}

func (d *shimDriver) Get(req Request) Response {
	var res Response
	v, err := d.d.Get(req.Name)
	if err != nil {
		res.Err = err.Error()
		return res
	}
	res.Volume = &Volume{
		Name:       v.Name(),
		Mountpoint: v.Path(),
	}
	return res
}

func (d *shimDriver) Remove(req Request) Response {
	var res Response
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

func (d *shimDriver) Path(req Request) Response {
	var res Response
	v, err := d.d.Get(req.Name)
	if err != nil {
		res.Err = err.Error()
		return res
	}
	res.Mountpoint = v.Path()
	return res
}

func (d *shimDriver) Mount(req Request) Response {
	var res Response
	v, err := d.d.Get(req.Name)
	if err != nil {
		res.Err = err.Error()
		return res
	}
	pth, err := v.Mount()
	if err != nil {
		res.Err = err.Error()
	}
	res.Mountpoint = pth
	return res
}

func (d *shimDriver) Unmount(req Request) Response {
	var res Response
	v, err := d.d.Get(req.Name)
	if err != nil {
		res.Err = err.Error()
		return res
	}
	if err := v.Unmount(); err != nil {
		res.Err = err.Error()
	}
	return res
}
