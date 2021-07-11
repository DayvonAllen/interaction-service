## Routes(All routes are protected)
- Send message:
    - `http://localhost:8083/messages/`
    - JSON:
```
{
    "to": <username>,
    "content": <message>
}
```    
- Get Conversation:
    - `http://localhost:8083/conversation/<username>`
- Get all conversation preview data:
    - `http://localhost:8083/conversation`
- Delete message by ID:
    - `http://localhost:8083/messages`
    - JSON:
```
{
    "messageId": "60ea310d029b10b21cef2447"
}
```   
- Delete messages by Ids:
    - `http://localhost:8083/messages/multi`
    - JSON:
```
{
    "ids": [
        {
            "messageId": "60ea73b64dc5537ff2823733"
        },
       {
            "messageId": "60ea73bd4dc5537ff2823737"
        },
        {
            "messageId": "60ea73bd4dc5537ff2823738"
        }
    ]
}
```    
---
