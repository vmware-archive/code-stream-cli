# Copyright 2019 VMware, Inc.
# SPDX-License-Identifier: 
FROM scratch
COPY cs-cli /
ENTRYPOINT [ "/cs-cli" ]