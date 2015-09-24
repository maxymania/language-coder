/*
      Copyright(C) 2015 Simon Schmidt

      This Source Code Form is subject to the
      terms of the Mozilla Public License, v.
      2.0. If a copy of the MPL was not
      distributed with this file, You can
      obtain one at
      http://mozilla.org/MPL/2.0/.
*/

package lowrepr


/*
 USER PROVIDES: A function, that represents an Operation.
 The function returns an Op object, that is implements the given
 operation. The .In and .Out field is ignored.
 */
type Operation func(v []ValueType) (Op,[]ValueType,error)

func op_move(v []ValueType) (Op,[]ValueType,error) {
	return Op{Code:OC_MOVE},v,nil
}

/*
 USER PROVIDES: An object that implements Operations and policies
 for implicit type casting.
 */
type Helper interface{
	/*
	 Performs a Typecheck and produces an error if the types don't match.
	 A non-nil Operation can be returned, if a implicit type-cast has to
	 be performed.
	 */
	Check(src, dest ValueType) (conv Operation,err error)
	Op(opcode int,u interface{}) (op Operation,err error)

	/*
	Handling of General Purpose Register types.
	*/
	GprType(used ValueType) ValueType
	IsGpr(gpr, used ValueType) bool
}



