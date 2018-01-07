resource "aws_vpc" "mongo" {
  cidr_block = "10.0.0.0/16"

  tags {
    Name = "${var.env}-mongo"
  }
}

resource "aws_internet_gateway" "default" {
  vpc_id = "${aws_vpc.mongo.id}"
}

resource "aws_route_table" "public" {
  vpc_id = "${aws_vpc.mongo.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.default.id}"
  }
}

resource "aws_subnet" "public" {
  count                   = 3
  vpc_id                  = "${aws_vpc.mongo.id}"
  cidr_block              = "${lookup(var.public_subnet_cidrs, format("zone%d", count.index))}"
  availability_zone       = "${lookup(var.zones, format("zone%d", count.index))}"
  depends_on              = ["aws_internet_gateway.default"]

  tags {
    Name = "${var.env}-public-${lookup(var.zones, format("zone%d", count.index))}"
  }
}

resource "aws_route_table_association" "public" {
  count          = 3
  subnet_id      = "${element(aws_subnet.public.*.id, count.index)}"
  route_table_id = "${aws_route_table.public.id}"
}

resource "aws_eip" "mongo" {
  count = 3
  vpc   = true
}

resource "aws_nat_gateway" "mongo" {
  count         = 3
  allocation_id = "${element(aws_eip.mongo.*.id, count.index)}"
  subnet_id     = "${element(aws_subnet.public.*.id, count.index)}"
}

resource "aws_subnet" "private" {
  count                   = 3
  vpc_id                  = "${aws_vpc.mongo.id}"
  cidr_block              = "${lookup(var.private_subnet_cidrs, format("zone%d", count.index))}"
  availability_zone       = "${lookup(var.zones, format("zone%d", count.index))}"
  depends_on              = ["aws_nat_gateway.mongo"]

  tags {
    Name = "${var.env}-private-${lookup(var.zones, format("zone%d", count.index))}"
  }
}

resource "aws_route_table" "private" {
  vpc_id = "${aws_vpc.mongo.id}"
  count  = 3
}

resource "aws_route" "private" {
  count                  = 3
  route_table_id         = "${element(aws_route_table.private.*.id, count.index)}"
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = "${element(aws_nat_gateway.mongo.*.id, count.index)}"
}

resource "aws_route_table_association" "private" {
  count          = 3
  subnet_id      = "${element(aws_subnet.private.*.id, count.index)}"
  route_table_id = "${element(aws_route_table.private.*.id, count.index)}"
}
