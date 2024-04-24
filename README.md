# waracle-test

### Environment Setup (Development)

1. Install Terraform
2. Install gcloud
3. Install Go 1.22+
4. Create gcloud configuration `gcloud config configurations create waracle-test-dev`
5. Set gcloud project `gcloud config set project waracle-test-dev`
6. Authenticate gcloud `gcloud auth login`

### New Environment Set up (GCP project)
1. Ensure correct gcloud project is configured
2. Enable services `./scripts/enable-services`
3. Manually create terraform state bucket matching naming convention `waracle-test-$STAGE-terraform`
   - Set region to 'europe-west2'
   - Storage class 'Standard'
   - Enforce public access prevention
   - Uniform access
   - Create


##### Dev deploy
`deploy --stage dev`

##### Prod deploy (NOT IMPLEMENTED)
`deploy --stage prod`
> You will be asked to type 'yes' after confirming resource changes