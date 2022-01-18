## DEVELOPMENT


### Authentication
The exporter uses [Application Default Credentials (ADCs)](https://cloud.google.com/docs/authentication/production#finding_credentials_automatically) to simplify authentication by finding credentials automatically.


On local development machine, you can authenticate using the following command, to establish your user account as ADCs
```
gcloud auth application-default login
```

On K8s, you need to create a Google service account, supply it using a secret and mount the secret to the exporter deployment.

```yaml
          env:
            - name: GOOGLE_APPLICATION_CREDENTIALS
              value: /secrets/gke-sa-client-secret.json
              
          volumeMounts:
            - name: secrets
              mountPath: /secrets              


       volumes:
        - name: secrets
          secret:
            secretName: gke-info-exporter-credentials
```            