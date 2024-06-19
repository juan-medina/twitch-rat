# Copyright (c) 2024 Juan Medina
#
#  All rights reserved. This software and related documentation are proprietary to Juan Medina.
#
#  This source code is for internal use only and may not be copied, modified, or distributed
#  without the express written permission of Juan Medina. Any use of this software for any
#  purpose other than its intended use is strictly prohibited and may result in severe civil
#  and criminal penalties.
#
#  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
#  INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
#  PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE
#  FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
#  OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
#  DEALINGS IN THE SOFTWARE.

if __name__ == "__main__":
    with open("internal/app/embed/version.txt", "r") as f:
        version = f.read()
        f.close()

        tokens = version.split(".")
        major = int(tokens[0])
        minor = int(tokens[1])
        release = int(tokens[2])
        patch = int(tokens[3])

        patch += 1

        version = f"{major}.{minor}.{release}.{patch}"

        with open("internal/app/embed/version.txt", "w") as f:
            f.write(version)
            f.close()
