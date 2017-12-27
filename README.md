#echo message center
##HTTP RESTful API
<pre>
URL: /publish?channel=channel1&channel=channel2
METHOD: POST
REQUEST:
  PARAMS:NONE
  BODY:
  ---------------Body Format-------------------
  	----------------------------------------------------------------------------------------------------------------------
	|					PROTOCOL HEADER		                                |   PROTOCOL BODY    |
	----------------------------------------------------------------------------------------------------------------------
	| FLAG | LENGTH | CHECKSUM | VERSION | COMMANDCODE | ERRORCODE | TEXTDATALENGTH | BINDATALENGTH | TEXTDATA | BINDATA |
	----------------------------------------------------------------------------------------------------------------------
	|  4B  |   4B   |    4B    |    4B   |     4B      |     4B    |       4B       |      4B       |  Unknown | Unknown |
	----------------------------------------------------------------------------------------------------------------------
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
----------------------------------------------------------------------------------------------------------------------
|					PROTOCOL HEADER		                                |   PROTOCOL BODY    |
----------------------------------------------------------------------------------------------------------------------
| FLAG | LENGTH | CHECKSUM | VERSION | COMMANDCODE | ERRORCODE | TEXTDATALENGTH | BINDATALENGTH | TEXTDATA | BINDATA |
----------------------------------------------------------------------------------------------------------------------
|  4B  |   4B   |    4B    |    4B   |     4B      |     4B    |       4B       |      4B       |  Unknown | Unknown |
----------------------------------------------------------------------------------------------------------------------
------------------------------------------------
</pre>
