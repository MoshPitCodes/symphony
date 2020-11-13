package apiserver

import (
	"context"
	"fmt"
	"net"

	"github.com/erkrnt/symphony/api"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServerAPIServer struct {
	APIServer *APIServer
}

func (s *GRPCServerAPIServer) GetClusters(ctx context.Context, in *api.RequestClusters) (*api.ResponseClusters, error) {
	clusters, err := s.APIServer.Resources.getClusters()

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	res := &api.ResponseClusters{
		Results: clusters,
	}

	return res, nil
}

func (s *GRPCServerAPIServer) GetLogicalVolume(ctx context.Context, in *api.RequestLogicalVolume) (*api.LogicalVolume, error) {
	lvID, err := uuid.Parse(in.ID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	lv, err := s.APIServer.Resources.getLogicalVolumeByID(lvID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	if lv == nil {
		st := status.New(codes.NotFound, "invalid_logical_volume_id")

		return nil, st.Err()
	}

	return lv, nil
}

func (s *GRPCServerAPIServer) GetLogicalVolumes(ctx context.Context, in *api.RequestLogicalVolumes) (*api.ResponseLogicalVolumes, error) {
	lvs, err := s.APIServer.Resources.getLogicalVolumes()

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	res := &api.ResponseLogicalVolumes{
		Results: lvs,
	}

	return res, nil
}

func (s *GRPCServerAPIServer) GetPhysicalVolume(ctx context.Context, in *api.RequestPhysicalVolume) (*api.PhysicalVolume, error) {
	pvID, err := uuid.Parse(in.ID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	pv, err := s.APIServer.Resources.getPhysicalVolumeByID(pvID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	if pv == nil {
		st := status.New(codes.NotFound, "invalid_physical_volume_id")

		return nil, st.Err()
	}

	return pv, nil
}

func (s *GRPCServerAPIServer) GetService(ctx context.Context, in *api.RequestService) (*api.Service, error) {
	serviceID, err := uuid.Parse(in.ServiceID)

	if err != nil {
		st := status.New(codes.InvalidArgument, "invalid_service_id")

		return nil, st.Err()
	}

	service, err := s.APIServer.Resources.getServiceByID(serviceID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	if service == nil {
		st := status.New(codes.NotFound, "invalid_service")

		return nil, st.Err()
	}

	return service, nil
}

func (s *GRPCServerAPIServer) GetServices(ctx context.Context, in *api.RequestServices) (*api.ResponseServices, error) {
	services, err := s.APIServer.Resources.getServices()

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	res := &api.ResponseServices{
		Results: services,
	}

	return res, nil
}

func (s *GRPCServerAPIServer) GetPhysicalVolumes(ctx context.Context, in *api.RequestPhysicalVolumes) (*api.ResponsePhysicalVolumes, error) {
	pvs, err := s.APIServer.Resources.getPhysicalVolumes()

	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	res := &api.ResponsePhysicalVolumes{
		Results: pvs,
	}

	return res, nil
}

func (s *GRPCServerAPIServer) GetVolumeGroup(ctx context.Context, in *api.RequestVolumeGroup) (*api.VolumeGroup, error) {
	vgID, err := uuid.Parse(in.ID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())
		return nil, st.Err()
	}

	vg, err := s.APIServer.Resources.getVolumeGroupByID(vgID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	if vg == nil {
		st := status.New(codes.NotFound, "invalid_volume_group_id")
		return nil, st.Err()
	}

	return vg, nil
}

func (s *GRPCServerAPIServer) GetVolumeGroups(ctx context.Context, in *api.RequestVolumeGroups) (*api.ResponseVolumeGroups, error) {
	vgs, err := s.APIServer.Resources.getVolumeGroups()

	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	res := &api.ResponseVolumeGroups{
		Results: vgs,
	}

	return res, nil
}

func (s *GRPCServerAPIServer) NewCluster(ctx context.Context, in *api.RequestNewCluster) (*api.Cluster, error) {
	clusterID := uuid.New()

	c := &api.Cluster{
		ID:     clusterID.String(),
		Status: api.ResourceStatus_CREATE_COMPLETED,
	}

	clusterSaveErr := s.APIServer.Resources.saveCluster(c)

	if clusterSaveErr != nil {
		st := status.New(codes.Internal, clusterSaveErr.Error())

		return nil, st.Err()
	}

	return c, nil
}

func (s *GRPCServerAPIServer) NewLogicalVolume(ctx context.Context, in *api.RequestNewLogicalVolume) (*api.LogicalVolume, error) {
	vgID, err := uuid.Parse(in.VolumeGroupID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	vg, err := s.APIServer.Resources.getVolumeGroupByID(vgID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	if vg == nil {
		st := status.New(codes.NotFound, "invalid_volume_group_id")

		return nil, st.Err()
	}

	pvID, err := uuid.Parse(vg.PhysicalVolumeID)

	if err != nil {
		st := status.New(codes.InvalidArgument, "invalid_service_id")

		return nil, st.Err()
	}

	pv, err := s.APIServer.Resources.getPhysicalVolumeByID(pvID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	if pv == nil {
		st := status.New(codes.NotFound, "invalid_physical_volume_id")

		return nil, st.Err()
	}

	lvID := uuid.New()

	lv := &api.LogicalVolume{
		ID:            lvID.String(),
		Size:          in.Size,
		Status:        api.ResourceStatus_REVIEW_IN_PROGRESS,
		VolumeGroupID: vg.ID,
	}

	saveErr := s.APIServer.Resources.saveLogicalVolume(lv)

	if saveErr != nil {
		return nil, saveErr
	}

	return lv, nil
}

func (s *GRPCServerAPIServer) NewPhysicalVolume(ctx context.Context, in *api.RequestNewPhysicalVolume) (*api.PhysicalVolume, error) {
	serviceID, err := uuid.Parse(in.ServiceID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())
		return nil, st.Err()
	}

	service, err := s.APIServer.Resources.getServiceByID(serviceID)

	if service == nil {
		st := status.New(codes.NotFound, "invalid_service_id")
		return nil, st.Err()
	}

	volumes, err := s.APIServer.Resources.getPhysicalVolumes()

	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	var volume *api.PhysicalVolume

	for _, v := range volumes {
		if in.DeviceName == v.DeviceName && service.ID == v.ServiceID {
			volume = v
		}
	}

	if volume != nil {
		st := status.New(codes.AlreadyExists, "physical_volume_exists")
		return nil, st.Err()
	}

	pvID := uuid.New()

	pv := &api.PhysicalVolume{
		DeviceName: in.DeviceName,
		ID:         pvID.String(),
		ServiceID:  service.ID,
	}

	saveErr := s.APIServer.Resources.savePhysicalVolume(pv)

	if saveErr != nil {
		return nil, saveErr
	}

	return pv, nil
}

func (s *GRPCServerAPIServer) NewService(ctx context.Context, in *api.RequestNewService) (*api.ResponseNewService, error) {
	clusterID, err := uuid.Parse(in.ClusterID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	cluster, err := s.APIServer.Resources.getClusterByID(clusterID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	if cluster == nil {
		st := status.New(codes.NotFound, "invalid_cluster")

		return nil, st.Err()
	}

	serviceID := uuid.New()

	service := &api.Service{
		ClusterID: cluster.ID,
		ID:        serviceID.String(),
		Status:    api.ResourceStatus_CREATE_COMPLETED,
		Type:      in.ServiceType,
	}

	serviceSaveErr := s.APIServer.Resources.saveService(service)

	if serviceSaveErr != nil {
		st := status.New(codes.Internal, serviceSaveErr.Error())

		return nil, st.Err()
	}

	res := &api.ResponseNewService{
		ClusterID: cluster.ID,
		ServiceID: service.ID,
	}

	if in.JoinAddr != "" {
		joinAddr, err := net.ResolveTCPAddr("tcp", in.JoinAddr)

		if err != nil {
			return nil, err
		}

		restartErr := s.APIServer.restart(joinAddr)

		if restartErr != nil {
			return nil, restartErr
		}

		return res, nil
	}

	local := s.APIServer.Node.Serf.Memberlist().LocalNode()

	res.JoinAddr = local.Addr.String()

	return res, nil
}

func (s *GRPCServerAPIServer) NewVolumeGroup(ctx context.Context, in *api.RequestNewVolumeGroup) (*api.VolumeGroup, error) {
	physicalVolumeID, err := uuid.Parse(in.PhysicalVolumeID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())
		return nil, st.Err()
	}

	physicalVolume, err := s.APIServer.Resources.getPhysicalVolumeByID(physicalVolumeID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	if physicalVolume == nil {
		st := status.New(codes.NotFound, "invalid_physical_volume_id")
		return nil, st.Err()
	}

	physicalVolumeServiceID, err := uuid.Parse(physicalVolume.ServiceID)

	if err != nil {
		st := status.New(codes.InvalidArgument, "invalid_service_id")
		return nil, st.Err()
	}

	service, err := s.APIServer.Resources.getServiceByID(physicalVolumeServiceID)

	if service == nil {
		st := status.New(codes.NotFound, "invalid_service_id")
		return nil, st.Err()
	}

	volumeGroupID := uuid.New()

	vg := &api.VolumeGroup{
		ID:               volumeGroupID.String(),
		PhysicalVolumeID: physicalVolume.ID,
		Status:           api.ResourceStatus_REVIEW_IN_PROGRESS,
	}

	saveErr := s.APIServer.Resources.saveVolumeGroup(vg)

	if saveErr != nil {
		return nil, saveErr
	}

	return vg, nil
}

func (s *GRPCServerAPIServer) RemoveLogicalVolume(ctx context.Context, in *api.RequestLogicalVolume) (*api.ResponseStatus, error) {
	logicalVolumeID, err := uuid.Parse(in.ID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())
		return nil, st.Err()
	}

	lv, err := s.APIServer.Resources.getLogicalVolumeByID(logicalVolumeID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	if lv == nil {
		st := status.New(codes.NotFound, "invalid_logical_volume_id")
		return nil, st.Err()
	}

	resourceKey := fmt.Sprintf("/LogicalVolume/%s", lv.ID)

	delErr := s.APIServer.Resources.removeResource(resourceKey)

	if delErr != nil {
		st := status.New(codes.Internal, delErr.Error())

		return nil, st.Err()
	}

	// TODO : emit an event change to block services

	res := &api.ResponseStatus{SUCCESS: true}

	return res, nil
}

func (s *GRPCServerAPIServer) RemovePhysicalVolume(ctx context.Context, in *api.RequestPhysicalVolume) (*api.ResponseStatus, error) {
	physicalVolumeID, err := uuid.Parse(in.ID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	pv, err := s.APIServer.Resources.getPhysicalVolumeByID(physicalVolumeID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	if pv == nil {
		st := status.New(codes.NotFound, "invalid_physical_volume_id")

		return nil, st.Err()
	}

	serviceID, err := uuid.Parse(pv.ServiceID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	service, err := s.APIServer.Resources.getServiceByID(serviceID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	if service == nil {
		st := status.New(codes.NotFound, "invalid_service_id")

		return nil, st.Err()
	}

	resourceKey := fmt.Sprintf("/PhysicalVolume/%s", pv.ID)

	delRes := s.APIServer.Resources.removeResource(resourceKey)

	if delRes != nil {
		st := status.New(codes.Internal, delRes.Error())

		return nil, st.Err()
	}

	res := &api.ResponseStatus{SUCCESS: true}

	return res, nil
}

func (s *GRPCServerAPIServer) RemoveService(ctx context.Context, in *api.RequestService) (*api.ResponseStatus, error) {
	clusterID, err := uuid.Parse(in.ClusterID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	cluster, err := s.APIServer.Resources.getClusterByID(clusterID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	if cluster == nil {
		st := status.New(codes.NotFound, "cluster_not_initialized")

		return nil, st.Err()
	}

	serviceID, err := uuid.Parse(in.ServiceID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	service, err := s.APIServer.Resources.getServiceByID(serviceID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	member, err := s.APIServer.Node.GetMember(service)

	if member != nil {
		st := status.New(codes.Unavailable, "service_unavailable")

		return nil, st.Err()
	}

	serviceAddr := member.Tags["ServiceAddr"]

	leaveAddr, err := net.ResolveTCPAddr("tcp", serviceAddr)

	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(leaveAddr.String(), grpc.WithInsecure())

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), GrpcContextTimeout)

	defer cancel()

	leaveOpts := &api.RequestServiceLeave{
		ServiceID: service.ID,
	}

	if service.Type == api.ServiceType_APISERVER {
		peer := api.NewAPIServerClient(conn)

		_, err := peer.ServiceLeave(ctx, leaveOpts)

		if err != nil {
			return nil, err
		}
	}

	if service.Type == api.ServiceType_BLOCK {
		peer := api.NewBlockClient(conn)

		_, err := peer.ServiceLeave(ctx, leaveOpts)

		if err != nil {
			return nil, err
		}
	}

	// TODO : replace with method on Resources

	resourceKey := fmt.Sprintf("/Service/%s", service.ID)

	delErr := s.APIServer.Resources.removeResource(resourceKey)

	if delErr != nil {
		st := status.New(codes.Internal, delErr.Error())

		return nil, st.Err()
	}

	res := &api.ResponseStatus{SUCCESS: true}

	return res, nil
}

func (s *GRPCServerAPIServer) RemoveVolumeGroup(ctx context.Context, in *api.RequestVolumeGroup) (*api.ResponseStatus, error) {
	volumeGroupID, err := uuid.Parse(in.ID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	vg, err := s.APIServer.Resources.getVolumeGroupByID(volumeGroupID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	if vg == nil {
		st := status.New(codes.NotFound, "invalid_volume_group_id")

		return nil, st.Err()
	}

	resourceKey := fmt.Sprintf("/VolumeGroup/%s", vg.ID)

	delErr := s.APIServer.Resources.removeResource(resourceKey)

	if delErr != nil {
		st := status.New(codes.Internal, delErr.Error())

		return nil, st.Err()
	}

	res := &api.ResponseStatus{SUCCESS: true}

	return res, nil
}

func (s *GRPCServerAPIServer) ServiceJoin(ctx context.Context, in *api.RequestServiceJoin) (*api.ResponseServiceJoin, error) {
	clusterID, err := uuid.Parse(in.ClusterID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	cluster, err := s.APIServer.Resources.getClusterByID(clusterID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	if cluster == nil {
		st := status.New(codes.NotFound, "invalid_cluster_id")

		return nil, st.Err()
	}

	serviceID, err := uuid.Parse(in.ServiceID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	service, err := s.APIServer.Resources.getServiceByID(serviceID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	if service == nil {
		st := status.New(codes.InvalidArgument, "invalid_service_id")

		return nil, st.Err()
	}

	local := s.APIServer.Node.Serf.Memberlist().LocalNode()

	localAddr := local.FullAddress().Addr

	res := &api.ResponseServiceJoin{
		PeerAddr: localAddr,
	}

	fields := logrus.Fields{
		"PeerAddr": localAddr,
	}

	logrus.WithFields(fields).Debug("ServiceJoin successful request for service")

	return res, nil
}

func (s *GRPCServerAPIServer) ServiceLeave(ctx context.Context, in *api.RequestServiceLeave) (*api.ResponseStatus, error) {
	clusterID, err := uuid.Parse(in.ClusterID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	cluster, err := s.APIServer.Resources.getClusterByID(clusterID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	if cluster == nil {
		st := status.New(codes.NotFound, "invalid_cluster_id")

		return nil, st.Err()
	}

	serviceID, err := uuid.Parse(in.ServiceID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	if serviceID != *s.APIServer.Node.Key.ServiceID {
		st := status.New(codes.PermissionDenied, err.Error())

		return nil, st.Err()
	}

	leaveErr := s.APIServer.Node.Serf.Leave()

	if leaveErr != nil {
		st := status.New(codes.Internal, leaveErr.Error())

		return nil, st.Err()
	}

	shutdownErr := s.APIServer.Node.Serf.Shutdown()

	if shutdownErr != nil {
		st := status.New(codes.Internal, shutdownErr.Error())

		return nil, st.Err()
	}

	logrus.Debug("Service 'apiserver' has left the cluster.")

	res := &api.ResponseStatus{SUCCESS: true}

	return res, nil
}
