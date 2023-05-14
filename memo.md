```
SQS_QUEUE_URL=....
SQS_MESSAGE="{"key1": "value1", "key2": "value2", "key3": "value3"}"

aws sqs send-message \
    --queue-url ${SQS_QUEUE_URL} \
    --message-body "${SQS_MESSAGE}"

for i in {1..15}; do
    aws sqs send-message \
        --queue-url ${SQS_QUEUE_URL} \
        --message-body "{\"key\": \"value $i\"}"
done


```