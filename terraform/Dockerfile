FROM hashicorp/terraform:0.11.14:light
COPY src /tf
WORKDIR /tf/s3-bucket
CMD [ "init" ]
CMD [ "plan" ]
CMD [ "apply", "-auto-approve"]
WORKDIR /tf/queue
CMD [ "/usr/bin/sleep 30 &&", "init" ]
CMD [ "apply", "-auto-approve"]
