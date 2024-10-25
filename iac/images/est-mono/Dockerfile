FROM openjdk:23-oracle
LABEL authors="vikazeni"

COPY ./est-mono/target/est-mono-0.0.1.jar /opt/app.jar

ENTRYPOINT ["java", "-jar", "/opt/app.jar"]
