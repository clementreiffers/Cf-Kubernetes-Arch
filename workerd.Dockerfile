FROM clementreiffers/worker-builder AS builder

COPY ./ ./

RUN workerd compile config.capnp > serv.out

FROM clementreiffers/worker-runner AS runner

COPY --from=builder serv.out .

CMD ["./serv.out"]
