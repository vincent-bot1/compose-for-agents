# ---- Stage 1: Build the app ----
FROM eclipse-temurin:17-jdk AS builder

WORKDIR /build

# Only copy metadata first to cache dependencies
COPY pom.xml ./
COPY .mvn .mvn
COPY mvnw ./
RUN ./mvnw dependency:go-offline -B

# Now copy the rest of the app
COPY src ./src
RUN ./mvnw package -DskipTests -B

# ---- Stage 2: Runtime image ----
FROM eclipse-temurin:17-jre AS app

WORKDIR /app

# Copy only the built jar from the builder stage
COPY --from=builder /build/target/*.jar app.jar

ENV SPRING_PROFILES_ACTIVE=prod
COPY <<EOF entrypoint.sh
#!/bin/sh
# spring.ai expects the base URL to omit the /v1/? suffix
export MODEL_RUNNER_URL=\${MODEL_RUNNER_URL%/}
export MODEL_RUNNER_URL=\${MODEL_RUNNER_URL%/v1}
exec java -jar app.jar
EOF
RUN chmod +x entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]
