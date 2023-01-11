FROM rust AS build
WORKDIR /opt/hopfencloud/

# Install rorm-cli
RUN cargo install rorm-cli --version 0.5.0 --root .

# For caching target/
RUN cargo init
COPY Cargo.toml ./Cargo.toml
RUN cargo build -r

# Copy project
COPY . ./

# Build project
RUN cargo build -r

FROM debian

# Copy rorm-cli
COPY --from=build /opt/hopfencloud/rorm-cli ./

# Copy binary
COPY --from=build /opt/hopfencloud/target/release/hopfencloud ./

ENTRYPOINT ["./hopfencloud"]

CMD ["start"]
