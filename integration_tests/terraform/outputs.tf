output "vpc_id" {
  value = "${aws_vpc.mongo.id}"
}

output "private_subnet_1" {
  value = "${aws_subnet.private.0.id}"
}

output "private_subnet_2" {
  value = "${aws_subnet.private.1.id}"
}

output "private_subnet_3" {
  value = "${aws_subnet.private.2.id}"
}
