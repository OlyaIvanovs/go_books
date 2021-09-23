FROM scratch
ADD go_books /go_books
CMD ["/go_books"]
EXPOSE 8080