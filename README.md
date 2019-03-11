# function.go

Simple Google Cloud Function that triggers on http and sends email using [Mailgun](https://www.mailgun.com/).

```
gcloud functions deploy mailgun \
        --env-vars-file env.yaml
        --entry-point Mail \
        --region europe-west1 \
        --runtime go111 \
        --trigger-http
```