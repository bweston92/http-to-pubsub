# HTTP to GCP PubSub Event Publisher

Relay HTTP requests to Google PubSub.

## Summary

This application consumes webhooks and will publish them to Google Cloud PubSub.

## Sending messages

```shell

```

## Running the server

### Docker image locally

Expected environment `GOOGLE_APPLICATION_CREDENTIALS` to point to service account file.

```shell
docker run -p 8000:8000 \
  -v "${GOOGLE_APPLICATION_CREDENTIALS}:/run/secrets/googleserviceaccount.json" \
  -e "GOOGLE_APPLICATION_CREDENTIALS=/run/secrets/googleserviceaccount.json" \
  ghcr.io/bweston92/http-to-pubsub:latest \
  -google-project-id "my-project" \
  -google-pubsub-topic "topic"
```

## Google Cloud PubSub topic and service account

You can use the Terraform module located in `terraform` folder.

```terraform
module "http_to_pubsub_module" {
  source = "git@github.com:bweston92/http-to-pubsub.git//terraform?ref=main"
  prefix = "userevents"
}
```

The module outputs include the Google Topic name which should be passed as `google-pubsub-topic`.

The module outputs include the Google Service Account key which should saved to a file and that file
should be set for `google-service-account`, alternatively you can use `GOOGLE_APPLICATION_CREDENTIALS`
environment variable to set the service account file location.

<details>
  <summary>Manual instructions</summary>
## 1. Create a Service Account
1. **Navigate to the IAM & Admin section**:
   - Open the Google Cloud Console.
   - Select "IAM & Admin" from the navigation menu.
   - Click on "Service Accounts".

2. **Create a new service account**:
    - Click "Create Service Account".
    - Enter `httptops-publisher` as the "Service account name".
    - Provide a description: "Publishes Ory events to PubSub".
    - Click "Create and continue".
    - Skip assigning roles for now and click "Done".

## 2. Create a Service Account Key
1. **Generate a new key for the service account**:
    - In the "Service Accounts" list, find the `httptops-publisher` service account you just created.
    - Click the three dots (menu) on the right side and select "Manage keys".
    - Click "Add Key" and then "Create new key".
    - Choose "JSON" as the key type.
    - Click "Create". A JSON file with the private key will be downloaded.

## 3. Create a Pub/Sub Topic
1. **Navigate to the Pub/Sub section**:
    - Open the Google Cloud Console.
    - Select "Pub/Sub" from the navigation menu.

2. **Create a new topic**:
    - Click "Create Topic".
    - Enter `httptops-events` as the topic name.
    - Click "Create".
    - After creating the topic, click on the topic name to access its settings.

3. **Set message retention duration**:
    - In the topic details page, click "Edit".
    - Set the "Message retention duration" to `86600s` (24 hours and 6 minutes).
    - Click "Save".

## 4. Assign Pub/Sub Publisher Role to the Service Account
1. **Assign IAM roles to the service account**:
    - In the Pub/Sub section, find and click on the `httptops-events` topic.
    - Click on the "Permissions" tab.
    - Click "Add Member".
    - In the "New members" field, enter the email of the `httptops-publisher` service account (it will look something like `httptops-publisher@<project-id>.iam.gserviceaccount.com`).
    - In the "Select a role" dropdown, choose "Pub/Sub Publisher".
    - Click "Save".

## 5. Outputs
1. **Service Account Key Output**:
    - The downloaded JSON file contains the private key and other credentials necessary for authenticating with the service account.

2. **Topic Output**:
    - The topic name `httptops-events` can be referenced directly within the Pub/Sub section.
</details>
