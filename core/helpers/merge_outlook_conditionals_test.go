package helpers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type outlookConditionalsTest struct {
	input  string
	output string
}

func TestMergeOutlookConditionals(t *testing.T) {
	var testValues = []outlookConditionalsTest{
		{
			input:  "<![endif]--><!--[if mso | IE]>",
			output: "",
		},
		{
			input: `
		</tr>
		<![endif]-->
		<!--[if mso | IE]>
		</td>`,
			output: "\n\t\t</tr>\n\t\t\n\t\t</td>",
		},
		{
			input:  "</div>\n              <!--[if mso | IE]>\n            </td>\n          <![endif]-->\n              <!--[if mso | IE]>\n        </tr>\n      <![endif]-->\n              <!--[if mso | IE]>\n                  </table>\n                <![endif]-->\n            </td>\n          </tr>\n        </tbody>\n      </table>\n    </div>\n    <!--[if mso | IE]>\n          </td>",
			output: "</div>\n              <!--[if mso | IE]>\n            </td>\n          \n        </tr>\n      \n                  </table>\n                <![endif]-->\n            </td>\n          </tr>\n        </tbody>\n      </table>\n    </div>\n    <!--[if mso | IE]>\n          </td>",
		},
	}

	for _, tt := range testValues {
		t.Run(fmt.Sprintf("should remove unnecessary conditionals %s", tt.input), func(t *testing.T) {
			assert.Equal(t, MergeOutlookConditionals(tt.input), tt.output)
		})
	}
}
