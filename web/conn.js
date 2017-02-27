

var connection = new autobahn.Connection({
         url: 'ws://localhost:8080/centrum',
         realm: 'realm1'
      });

connection.onopen = function (session) {
  console.log("OPEN!!!")

   // 1) subscribe to a topic
   function onevent(args) {
      console.log("Event:", args[0]);
   }

   /*
   session.subscribe('com.myapp.hello', onevent);

   // 2) publish an event
   session.publish('com.myapp.hello', ['Hello, world!']);

   // 3) register a procedure for remoting
   function add2(args) {
      return args[0] + args[1];
   }
   session.register('com.myapp.add2', add2);

   // 4) call a remote procedure
   */

   console.log("XXX!")

   session.call('get_config').then(

     function (config_in) {
        console.log("Result:", config_in);
        obj = JSON.parse(config_in);
        connection.update_config(obj)
     }
   );
};



/*

function connect() {

var ws = new WebSocket("ws://localhost:8080/control");

ws.onopen = function()
{
  ws.send(JSON.stringify({
    "type":"GET_CONFIG"
  }));
  console.log("request sent")
  //alert("Message is sent...");
};

ws.onmessage = function (evt)
{
  var received_msg = evt.data;
  console.log("1")
  console.log(">> got something ::: "+evt.data)
  //alert("Message is received...");
};

ws.onclose = function()
{
  // websocket is closed.
  //alert("Connection is closed...");
  //console.log("connection closed")

  console.log('Socket is closed. Reconnect will be attempted in 1 second.');
  setTimeout(function() {
    connect();
  }, 1000)

};

ws.onerror = function(err) {
  console.error('Socket encountered error: ', err.message, 'Closing socket')
  ws.close()
};

}

connect()
*/
