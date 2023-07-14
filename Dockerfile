FROM --platform=linux/amd64 ubuntu
WORKDIR /
RUN apt update && apt-get install -y ca-certificates  
COPY ${PWD}/output/gcp-azure-demo /gcp-azure-demo 
ENTRYPOINT [ "./gcp-azure-demo" ]