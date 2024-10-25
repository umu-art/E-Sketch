FROM openjdk:23-oracle
LABEL authors="vikazeni"

COPY ./est-mono/target/est-mono-0.0.1.jar /opt/app.jar
COPY ./iac/images/est-mono/elastic-apm-agent-1.52.0.jar /opt/apm.jar

ENTRYPOINT ["java", "-javaagent:/opt/apm.jar", "-jar", "/opt/app.jar"]
