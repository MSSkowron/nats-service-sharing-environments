variable "stream_name" {
    description = "name of the stream"
    type = string
    default = "foo"
}

variable "stream_subjects" {
    description = "subjects of the stream, list"
    type = list(string)
    default = ["*"]
}

variable "servers" {
    description = "list of servers to connect to in a comma seperated list"
    type = string
    default = ""
}

variable "login" {
    description = "login"
    type = string
    default = "admin"
}

variable "password" {
    description = "password"
    type = string
    default = "admin"
}

