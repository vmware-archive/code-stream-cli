# Copyright 2021 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause
FROM scratch
COPY cs-cli /
ENTRYPOINT [ "/cs-cli" ]