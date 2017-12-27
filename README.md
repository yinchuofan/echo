# echo
基于websocket的消息订阅发布中心

#echo message center
##HTTP RESTful API
<pre>
URL: /publish
METHOD: POST
REQUEST:
  PARAMS:NONE
  BODY:
  ---------------Body Format-------------------
  channels : channel1,channel2,... \r\n
  message : json
    {
      "code" : string //message code
      ,"content" : json //message content
      ,"publisher" : string //message publisher client id
      ,"time" : datetime string //message send time, format '2016-01-01 10:10:10'
    }
  ------------------------------------------------
RESPONSE:
  BODY:
		"success" or "error"
</pre>

<pre>
URL: /online
METHOD: GET
REQUEST:
  PARAMS:{"clientid":string list} //client id list
  BODY:NONE
RESPONSE:
  BODY:
    "error" or online id list		
</pre>

##Websocket private protocol base on json
<pre>
Register to message center:
---------------Message Format-------------------
command : register \r\n
\r\n
clientid : string
------------------------------------------------
</pre>

<pre>
Subscribe message on channels:
---------------Message Format-------------------
command : subscribe \r\n
\r\n
channels : channel1,channel2,...
------------------------------------------------
</pre>

<pre>
Unsubscribe message on channels:
---------------Message Format-------------------
command : unsubscribe \r\n
\r\n
channels : channel1,channel2,...
------------------------------------------------
</pre>

<pre>
Publish message to channels:
---------------Message Format-------------------
command : publish \r\n
\r\n
channels : channel1,channel2,... \r\n
message : json
  {
    "code" : string //message code
    ,"content" : json //message content
    ,"publisher" : string //message publisher client id
    ,"time" : datetime string //message send time, format '2016-01-01 10:10:10'
  }
------------------------------------------------
</pre>
