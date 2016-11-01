var zmq = require('zmq')
  , sock = zmq.socket('sub');

url = 'tcp://0.0.0.0:8341'
sock.connect(url);
sock.subscribe('');
console.log('Subscriber connected to ' + url);

sock.on('connect', function(fd, ep) {console.log('connect, endpoint:', ep);});

sock.on('message', function(topic, message) {
	m=message.toString("utf-8")	
	console.log('[', topic.toString("utf-8"), ']:', message.toString("utf-8"));
	/*obj=JSON.parse(m)
	if (obj.hasOwnProperty("Rejection")) {
		payload = obj["Rejection"]["tx"]["payload"];
		console.log('------payload:', payload.toString("utf-8"));		
	}*/
});