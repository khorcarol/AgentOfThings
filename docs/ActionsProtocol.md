# Actions Protocol

## Overview

### Messages

All messages will be sent in the JSON format.
There are 6 types of messages that can be sent:

- **Functionality Query (FQ):** A request to ask for a nodes functionality,
this has the same format for every request.
- **Functionality Response (FR):**
A response to an FQ, contains data on the node type,
and information about what requests can be made, this includes:
  - An action ID
  - Action name
  - A description of the request
  - The type of action
  - Parameter details - description, data type
  - Return data details - description, data type
  - The nodes current state.
- **Action Call (AC):** This requests an action from a node, it will contain:
  - The action ID being called
  - A call ID
  - Any data the action needs
- **Action Response (AR):**
This is sent in once an action has completed, it will contain:
  - The call ID
  - Whether the action was successful
  - Any return data
  - Any error messages
- **Update Query (NU):** Requests the state of a node
- **Update Response (NU):** Responds to an update query with the state of the node

### Action Types

Message types give information about what the action should do.
Each type has expected parameters and returns that need to be implemented,

- Text: send a message
- Toggle: A toggle switch (with state)
- Integer range: Holds an integer value in a range
- Double range: Holds a double value in a range
- Custom: Does not have any expected parameters or returns

## Specification

### Messages

#### Functionality Query

```json
"mtype" = "fq",
```

#### Functionality Response

```JSON
"mtype": "fr",
"id": 1234,
"name": "Node Name",
"desc": "Node desc",
"actions": [
 {
  "id": 4321,
  "name": "Action Name",
  "desc": "A description of the action",
  "type": "data",
  "parameters": [
   {"id": 0, "name": "Parameter", "desc": "A parameter description", "type": "bool"},],
  "return": [{"name": "Return Value", desc: "", type: "double"},],
  "state": [
   {"id": 0, "name": "Action State", desc: "State specific to an action", type: true},
  ],
 },
 ...
],
"node_state": [
 {
  "id": 0,
  "name": "State name",
  "desc": "State description",
  "type": "double"
 }
]
```

#### Action Call

```JSON
"mtype": "ac",
"action_id": 4321,
"call_id": 6789,
"parameters": [
 {
  "id": 0,
  "data": 12,
 },
],
```

#### Action Response

```JSON
"mtype": "ar",
"action_id": 4321,
"call_id": 6789,
"success": 1,
"errors": [],
```

A "success" of 0 means that the action was successful,
any other number means that the action was not successful.

#### Update Query

```JSON
"mtype": "uq"
```

#### Update Response

```JSON
"mtype": "ur"
"node_state": []
"action_state": [
 {
  "id": 4321,
  "state": {
   "id": 0,
   "data": true
  }
 }, ...
],
"node_state": [
 {
  "id": 0,
  "data": 54.7,
 }
]
```

### Data types

There will be 5 different accepted data types:

- Integers - "int"
- Booleans - "bool"
- Double - "double"
- String - "str"
- Binary - "bin"

### Action Types

#### Message

Identifier: "message"

Required parameters:

- String: Message to be sent

Return values: None

State: None

#### Toggle

Identifier: "toggle"
Required parameters:

- Boolean: on/off

Return values: None
State:

- Boolean: toggle on or off

#### Integer Range

Description: Holds an integer whose max and min values are determined.
New value parameter must be checked server side to see if it is in the max range.

Identifier: "irange"

Required parameters:

- Integer: new value
Return values: None
State:
- Integer: Current value
- Integer: Min value
- Integer: Max value

#### Double Range

Description: Holds an double whose max and min values are determined.
New value parameter must be checked server side to see if it is in the max range.
Identifier: "drange"
Required parameters:

- Double: new value
Return values: None
State:
- Double: Current value
- Double: Min value
- Double: Max value
