package driver

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	"golang.org/x/net/context"
)

func (d *driver) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	resp := &csi.NodeGetInfoResponse{
		NodeId: d.config.NodeID,
	}
	return resp, nil
}

func (d *driver) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	return &csi.NodeStageVolumeResponse{}, nil
}

func (d *driver) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (d *driver) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	glog.V(5).Infof("NodePublishVolume: called with args %+v", *req)
	var err error
	rpath := "/"
	vfsOpt := make(map[string]any)
	mountOpt := make(map[string]any)

	//var b []byte

	if v, ok := req.VolumeContext["path"]; ok {
		rpath = v
	}
	if v, ok := req.VolumeContext["vfs"]; ok {
		if err = json.Unmarshal([]byte(v), &vfsOpt); err != nil {
			goto clean
		}
	}
	if req.Readonly {
		vfsOpt["ReadOnly"] = true
	}
	if v, ok := req.VolumeContext["mount"]; ok {
		if err = json.Unmarshal([]byte(v), &mountOpt); err != nil {
			goto clean
		}
	}
	glog.V(5).Infof("yes")
	if v, ok := req.VolumeCapability.GetAccessType().(*csi.VolumeCapability_Mount); ok {
		glog.V(5).Infof("YO: %+v", v.Mount.GetMountFlags())
		flags := v.Mount.GetMountFlags()

		s := strings.Split(strings.Trim(flags[0], " "), "--")
		for _, f := range s {
			if f == "" {
				continue
			}
			//glog.V(5).Infof("f %+v", f)
			f = strings.Trim(f, " ") 
			//glog.V(5).Infof("f %+v", f)
			k := strings.Split(f, " ")
			//glog.V(5).Infof("k %+v", k)
			b := strings.Split(k[0], "-")
			//glog.V(5).Infof("b %+v", b)
			name := ""
			for _, x := range b {
				name += strings.Title(x)
			}
			//glog.V(5).Infof("name %+v", name)
			
			if len(k) == 1 {
				mountOpt[name] = true
			} else {
				//glog.V(5).Infof("value %+v", k[1])
				mountOpt[name] = k[1]
			}
			//--allow-non-empty --allow-other --vfs-cache-mode full --dir-cache-time 10s
		}
		// for i, f := range flags {
		// 	//s := strings.Split(f, ":");
		// 	// if len(s) == 1 {
		// 	// 	continue
		// 	// }
		// 	mountOpt[f] =
		// 	// f.
		// 	// if !hasMountOption(mountOptions, f) {
		// 	// 	mountOptions = append(mountOptions, f)
		// 	// }
		// }
		//v.Mount.GetMountFlags()

		//glog.Infof("Uhu: %+v", v.Mount["mount_flags"]);
		// if err = v.Mount.XXX_Unmarshal(b); err != nil {
		// if err = json.Unmarshal([v.Mount.MountFlags, &mountOpt); err != nil {
		// 	glog.Infof("ERROR");
		// 	goto clean
		// }
	}

	//glog.V(5).Infof("YO: %+v", b);
	glog.V(5).Infof("YO: %+v", mountOpt)

	// if v, ok := req.VolumeContext["mount"]; ok {
	// 	if err = json.Unmarshal([]byte(v), &mountOpt); err != nil {
	// 		goto clean
	// 	}
	// }
	// test := req.VolumeCapability.GetMount();
	// glog.V(5).Infof("nount: %+v", test);
	// if err = json.Unmarshal([]byte(req.VolumeCapability.GetMount()), &mountOpt); err != nil {
	// 	goto clean
	// }

	if _, err = d.remoteCreate(req.VolumeId, req.VolumeContext); err != nil {
		goto clean
	}
	_, err = d.remoteMount(req.VolumeId, rpath, req.TargetPath, vfsOpt, mountOpt)
clean:
	glog.V(5).Infof("publish volume: %+v", err)
	return &csi.NodePublishVolumeResponse{}, err
}

func (d *driver) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	_, err := d.remoteUmount(req.TargetPath)
	glog.V(5).Infof("unpublish volume: %+v", err)
	return &csi.NodeUnpublishVolumeResponse{}, err
}

func (d *driver) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	cl := []csi.NodeServiceCapability_RPC_Type{}
	var caps []*csi.NodeServiceCapability
	for _, c := range cl {
		caps = append(caps, &csi.NodeServiceCapability{
			Type: &csi.NodeServiceCapability_Rpc{
				Rpc: &csi.NodeServiceCapability_RPC{
					Type: c,
				},
			},
		})
	}
	return &csi.NodeGetCapabilitiesResponse{Capabilities: caps}, nil
}

func (d *driver) NodeGetVolumeStats(ctx context.Context, in *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	return nil, fmt.Errorf("unsupported")
}

// NodeExpandVolume is only implemented so the driver can be used for e2e testing.
func (hp *driver) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	return nil, fmt.Errorf("unsupported")
}
