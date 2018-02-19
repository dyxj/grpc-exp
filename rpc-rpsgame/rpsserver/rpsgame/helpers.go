package rpsgame

import "github.com/sirupsen/logrus"

func getRequest(stream RpsSvc_GameServer) (*Req, error) {
	req, err := stream.Recv()
	if err != nil {
		logrus.Errorf("getRequest(), Error receiving: %v", err)
		return nil, err
	}
	logrus.Infof("getRequest(), Received: %v", req)
	return req, nil
}

func sendResponse(stream RpsSvc_GameServer, r *Resp) error {
	logrus.Infof("sendResponse(), Sending response: %v", r)
	err := stream.Send(r)

	if err != nil {
		logrus.Errorf("sendResponse(), Error sending: %v", err)
	}

	return err
}

func sendSign(stream RpsSvc_GameServer, sign Sign) error {
	rsign := &Resp_Sign{Sign: sign}
	resp := &Resp{Event: rsign}
	return sendResponse(stream, resp)
}

func sendState(stream RpsSvc_GameServer, state Resp_State) error {
	rstate := &Resp_Gstate{Gstate: state}
	resp := &Resp{Event: rstate}
	return sendResponse(stream, resp)
}
