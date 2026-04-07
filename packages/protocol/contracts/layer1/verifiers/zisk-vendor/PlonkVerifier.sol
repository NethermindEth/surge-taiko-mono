// SPDX-License-Identifier: GPL-3.0
/*
    Copyright 2021 0KIMS association.

    This file is generated with [snarkJS](https://github.com/iden3/snarkjs).

    snarkJS is a free software: you can redistribute it and/or modify it
    under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    snarkJS is distributed in the hope that it will be useful, but WITHOUT
    ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
    or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public
    License for more details.

    You should have received a copy of the GNU General Public License
    along with snarkJS. If not, see <https://www.gnu.org/licenses/>.
*/

pragma solidity >=0.7.0 <0.9.0;

contract PlonkVerifier {
    // Omega
    uint256 constant w1 =
        5_709_868_443_893_258_075_976_348_696_661_355_716_898_495_876_243_883_251_619_397_131_511_003_808_859;
    // Scalar field size
    uint256 constant q =
        21_888_242_871_839_275_222_246_405_745_257_275_088_548_364_400_416_034_343_698_204_186_575_808_495_617;
    // Base field size
    uint256 constant qf =
        21_888_242_871_839_275_222_246_405_745_257_275_088_696_311_157_297_823_662_689_037_894_645_226_208_583;

    // [1]_1
    uint256 constant G1x = 1;
    uint256 constant G1y = 2;
    // [1]_2
    uint256 constant G2x1 =
        10_857_046_999_023_057_135_944_570_762_232_829_481_370_756_359_578_518_086_990_519_993_285_655_852_781;
    uint256 constant G2x2 =
        11_559_732_032_986_387_107_991_004_021_392_285_783_925_812_861_821_192_530_917_403_151_452_391_805_634;
    uint256 constant G2y1 =
        8_495_653_923_123_431_417_604_973_247_489_272_438_418_190_587_263_600_148_770_280_649_306_958_101_930;
    uint256 constant G2y2 =
        4_082_367_875_863_433_681_332_203_403_145_435_568_316_851_327_593_401_208_105_741_076_214_120_093_531;

    // Verification Key data
    uint32 constant n = 16_777_216;
    uint16 constant nPublic = 1;
    uint16 constant nLagrange = 1;

    uint256 constant Qmx =
        15_728_078_222_621_428_176_160_834_771_923_049_445_267_207_177_815_563_086_756_963_442_982_537_322_555;
    uint256 constant Qmy =
        10_368_134_109_772_487_360_967_338_421_067_531_937_191_459_049_728_961_159_479_449_769_737_397_414_266;
    uint256 constant Qlx =
        6_266_903_827_059_262_482_682_470_537_784_423_436_414_327_073_503_210_506_957_823_438_202_715_076_858;
    uint256 constant Qly =
        9_930_975_912_673_080_337_533_028_672_610_162_866_366_123_480_191_667_478_375_977_723_280_338_682_627;
    uint256 constant Qrx =
        7_123_060_314_877_137_702_038_864_858_686_611_023_352_874_626_900_523_603_313_267_911_196_395_216_007;
    uint256 constant Qry =
        1_770_043_357_969_798_948_108_297_406_468_647_422_219_312_923_587_004_610_833_401_748_122_605_130_485;
    uint256 constant Qox =
        429_670_666_919_728_039_002_263_956_209_270_185_134_416_998_857_239_332_019_860_571_923_582_607_053;
    uint256 constant Qoy =
        11_024_111_753_296_480_975_480_562_321_760_673_479_457_513_637_396_733_200_807_061_983_645_283_791_758;
    uint256 constant Qcx =
        18_759_539_946_247_966_109_746_686_045_765_624_942_782_666_785_179_752_870_440_507_130_849_473_521_082;
    uint256 constant Qcy =
        3_887_166_330_695_441_169_968_198_076_272_085_557_699_876_114_516_022_583_985_382_505_001_285_608_090;
    uint256 constant S1x =
        14_529_212_361_763_788_633_811_869_842_580_421_415_134_713_687_938_570_863_064_659_590_110_934_903_740;
    uint256 constant S1y =
        15_130_539_367_607_292_893_898_013_497_807_485_005_818_775_366_682_240_859_886_811_844_322_272_868_379;
    uint256 constant S2x =
        88_575_883_471_378_427_909_504_541_476_422_256_478_241_201_143_920_420_405_836_063_934_407_519_240;
    uint256 constant S2y =
        3_538_048_320_551_462_152_620_488_185_558_124_425_114_005_176_021_275_627_227_409_411_627_962_175_942;
    uint256 constant S3x =
        18_339_485_453_472_040_659_646_089_143_258_381_971_176_107_908_300_548_337_940_061_504_006_306_000_991;
    uint256 constant S3y =
        401_317_798_322_731_345_836_300_355_016_089_251_665_269_411_090_279_566_075_041_830_406_289_046_834;
    uint256 constant k1 = 2;
    uint256 constant k2 = 3;
    uint256 constant X2x1 =
        21_831_381_940_315_734_285_607_113_342_023_901_060_522_397_560_371_972_897_001_948_545_212_302_161_822;
    uint256 constant X2x2 =
        17_231_025_384_763_736_816_414_546_592_865_244_497_437_017_442_647_097_510_447_326_538_965_263_639_101;
    uint256 constant X2y1 =
        2_388_026_358_213_174_446_665_280_700_919_698_872_609_886_601_280_537_296_205_114_254_867_301_080_648;
    uint256 constant X2y2 =
        11_507_326_595_632_554_467_052_522_095_592_665_270_651_932_854_513_688_777_769_618_397_986_436_103_170;

    // Proof calldata
    // Byte offset of every parameter of the calldata
    // Polynomial commitments
    uint16 constant pA = 4 + 0;
    uint16 constant pB = 4 + 64;
    uint16 constant pC = 4 + 128;
    uint16 constant pZ = 4 + 192;
    uint16 constant pT1 = 4 + 256;
    uint16 constant pT2 = 4 + 320;
    uint16 constant pT3 = 4 + 384;
    uint16 constant pWxi = 4 + 448;
    uint16 constant pWxiw = 4 + 512;
    // Opening evaluations
    uint16 constant pEval_a = 4 + 576;
    uint16 constant pEval_b = 4 + 608;
    uint16 constant pEval_c = 4 + 640;
    uint16 constant pEval_s1 = 4 + 672;
    uint16 constant pEval_s2 = 4 + 704;
    uint16 constant pEval_zw = 4 + 736;

    // Memory data
    // Challenges
    uint16 constant pAlpha = 0;
    uint16 constant pBeta = 32;
    uint16 constant pGamma = 64;
    uint16 constant pXi = 96;
    uint16 constant pXin = 128;
    uint16 constant pBetaXi = 160;
    uint16 constant pV1 = 192;
    uint16 constant pV2 = 224;
    uint16 constant pV3 = 256;
    uint16 constant pV4 = 288;
    uint16 constant pV5 = 320;
    uint16 constant pU = 352;

    uint16 constant pPI = 384;
    uint16 constant pEval_r0 = 416;
    uint16 constant pD = 448;
    uint16 constant pF = 512;
    uint16 constant pE = 576;
    uint16 constant pTmp = 640;
    uint16 constant pAlpha2 = 704;
    uint16 constant pZh = 736;
    uint16 constant pZhInv = 768;

    uint16 constant pEval_l1 = 800;

    uint16 constant lastMem = 832;

    function verifyProof(
        uint256[24] calldata _proof,
        uint256[1] calldata _pubSignals
    )
        public
        view
        returns (bool)
    {
        assembly {
            /////////
            // Computes the inverse using the extended euclidean algorithm
            /////////
            function inverse(a, q) -> inv {
                let t := 0
                let newt := 1
                let r := q
                let newr := a
                let quotient
                let aux

                for { } newr { } {
                    quotient := sdiv(r, newr)
                    aux := sub(t, mul(quotient, newt))
                    t := newt
                    newt := aux

                    aux := sub(r, mul(quotient, newr))
                    r := newr
                    newr := aux
                }

                if gt(r, 1) { revert(0, 0) }
                if slt(t, 0) { t := add(t, q) }

                inv := t
            }

            ///////
            // Computes the inverse of an array of values
            // See https://vitalik.ca/general/2018/07/21/starks_part_3.html in section where explain fields operations
            //////
            function inverseArray(pVals, n) {
                let pAux := mload(0x40) // Point to the next free position
                let pIn := pVals
                let lastPIn := add(pVals, mul(n, 32)) // Read n elements
                let acc := mload(pIn) // Read the first element
                pIn := add(pIn, 32) // Point to the second element
                let inv

                for { } lt(pIn, lastPIn) {
                    pAux := add(pAux, 32)
                    pIn := add(pIn, 32)
                } {
                    mstore(pAux, acc)
                    acc := mulmod(acc, mload(pIn), q)
                }
                acc := inverse(acc, q)

                // At this point pAux pint to the next free position we subtract 1 to point to the last used
                pAux := sub(pAux, 32)
                // pIn points to the n+1 element, we subtract to point to n
                pIn := sub(pIn, 32)
                lastPIn := pVals // We don't process the first element
                for { } gt(pIn, lastPIn) {
                    pAux := sub(pAux, 32)
                    pIn := sub(pIn, 32)
                } {
                    inv := mulmod(acc, mload(pAux), q)
                    acc := mulmod(acc, mload(pIn), q)
                    mstore(pIn, inv)
                }
                // pIn points to first element, we just set it.
                mstore(pIn, acc)
            }

            function checkField(v) {
                if iszero(lt(v, q)) {
                    mstore(0, 0)
                    return(0, 0x20)
                }
            }

            function checkPointBelongsToBN128Curve(p) {
                let x := calldataload(p)
                let y := calldataload(add(p, 32))

                // Check that the point is on the curve
                // y^2 = x^3 + 3
                let x3_3 := addmod(mulmod(x, mulmod(x, x, qf), qf), 3, qf)
                let y2 := mulmod(y, y, qf)

                if iszero(eq(x3_3, y2)) {
                    mstore(0, 0)
                    return(0, 0x20)
                }
            }

            function checkProofData() {
                // Check proof commitments belong to the bn128 curve
                checkPointBelongsToBN128Curve(pA)
                checkPointBelongsToBN128Curve(pB)
                checkPointBelongsToBN128Curve(pC)
                checkPointBelongsToBN128Curve(pZ)
                checkPointBelongsToBN128Curve(pT1)
                checkPointBelongsToBN128Curve(pT2)
                checkPointBelongsToBN128Curve(pT3)
                checkPointBelongsToBN128Curve(pWxi)
                checkPointBelongsToBN128Curve(pWxiw)

                // Check proof commitments coordinates are in the field
                checkField(calldataload(pA))
                checkField(calldataload(add(pA, 32)))
                checkField(calldataload(pB))
                checkField(calldataload(add(pB, 32)))
                checkField(calldataload(pC))
                checkField(calldataload(add(pC, 32)))
                checkField(calldataload(pZ))
                checkField(calldataload(add(pZ, 32)))
                checkField(calldataload(pT1))
                checkField(calldataload(add(pT1, 32)))
                checkField(calldataload(pT2))
                checkField(calldataload(add(pT2, 32)))
                checkField(calldataload(pT3))
                checkField(calldataload(add(pT3, 32)))
                checkField(calldataload(pWxi))
                checkField(calldataload(add(pWxi, 32)))
                checkField(calldataload(pWxiw))
                checkField(calldataload(add(pWxiw, 32)))

                // Check proof evaluations are in the field
                checkField(calldataload(pEval_a))
                checkField(calldataload(pEval_b))
                checkField(calldataload(pEval_c))
                checkField(calldataload(pEval_s1))
                checkField(calldataload(pEval_s2))
                checkField(calldataload(pEval_zw))
            }

            function calculateChallenges(pMem, pPublic) {
                let beta
                let aux

                let mIn := mload(0x40) // Pointer to the next free memory position

                // Compute challenge.beta & challenge.gamma
                mstore(mIn, Qmx)
                mstore(add(mIn, 32), Qmy)
                mstore(add(mIn, 64), Qlx)
                mstore(add(mIn, 96), Qly)
                mstore(add(mIn, 128), Qrx)
                mstore(add(mIn, 160), Qry)
                mstore(add(mIn, 192), Qox)
                mstore(add(mIn, 224), Qoy)
                mstore(add(mIn, 256), Qcx)
                mstore(add(mIn, 288), Qcy)
                mstore(add(mIn, 320), S1x)
                mstore(add(mIn, 352), S1y)
                mstore(add(mIn, 384), S2x)
                mstore(add(mIn, 416), S2y)
                mstore(add(mIn, 448), S3x)
                mstore(add(mIn, 480), S3y)

                mstore(add(mIn, 512), calldataload(add(pPublic, 0)))

                mstore(add(mIn, 544), calldataload(pA))
                mstore(add(mIn, 576), calldataload(add(pA, 32)))
                mstore(add(mIn, 608), calldataload(pB))
                mstore(add(mIn, 640), calldataload(add(pB, 32)))
                mstore(add(mIn, 672), calldataload(pC))
                mstore(add(mIn, 704), calldataload(add(pC, 32)))

                beta := mod(keccak256(mIn, 736), q)
                mstore(add(pMem, pBeta), beta)

                // challenges.gamma
                mstore(add(pMem, pGamma), mod(keccak256(add(pMem, pBeta), 32), q))

                // challenges.alpha
                mstore(mIn, mload(add(pMem, pBeta)))
                mstore(add(mIn, 32), mload(add(pMem, pGamma)))
                mstore(add(mIn, 64), calldataload(pZ))
                mstore(add(mIn, 96), calldataload(add(pZ, 32)))

                aux := mod(keccak256(mIn, 128), q)
                mstore(add(pMem, pAlpha), aux)
                mstore(add(pMem, pAlpha2), mulmod(aux, aux, q))

                // challenges.xi
                mstore(mIn, aux)
                mstore(add(mIn, 32), calldataload(pT1))
                mstore(add(mIn, 64), calldataload(add(pT1, 32)))
                mstore(add(mIn, 96), calldataload(pT2))
                mstore(add(mIn, 128), calldataload(add(pT2, 32)))
                mstore(add(mIn, 160), calldataload(pT3))
                mstore(add(mIn, 192), calldataload(add(pT3, 32)))

                aux := mod(keccak256(mIn, 224), q)
                mstore(add(pMem, pXi), aux)

                // challenges.v
                mstore(mIn, aux)
                mstore(add(mIn, 32), calldataload(pEval_a))
                mstore(add(mIn, 64), calldataload(pEval_b))
                mstore(add(mIn, 96), calldataload(pEval_c))
                mstore(add(mIn, 128), calldataload(pEval_s1))
                mstore(add(mIn, 160), calldataload(pEval_s2))
                mstore(add(mIn, 192), calldataload(pEval_zw))

                let v1 := mod(keccak256(mIn, 224), q)
                mstore(add(pMem, pV1), v1)

                // challenges.beta * challenges.xi
                mstore(add(pMem, pBetaXi), mulmod(beta, aux, q))

                // challenges.xi^n

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                aux := mulmod(aux, aux, q)

                mstore(add(pMem, pXin), aux)

                // Zh
                aux := mod(add(sub(aux, 1), q), q)
                mstore(add(pMem, pZh), aux)
                mstore(add(pMem, pZhInv), aux) // We will invert later together with lagrange pols

                // challenges.v^2, challenges.v^3, challenges.v^4, challenges.v^5
                aux := mulmod(v1, v1, q)
                mstore(add(pMem, pV2), aux)
                aux := mulmod(aux, v1, q)
                mstore(add(pMem, pV3), aux)
                aux := mulmod(aux, v1, q)
                mstore(add(pMem, pV4), aux)
                aux := mulmod(aux, v1, q)
                mstore(add(pMem, pV5), aux)

                // challenges.u
                mstore(mIn, calldataload(pWxi))
                mstore(add(mIn, 32), calldataload(add(pWxi, 32)))
                mstore(add(mIn, 64), calldataload(pWxiw))
                mstore(add(mIn, 96), calldataload(add(pWxiw, 32)))

                mstore(add(pMem, pU), mod(keccak256(mIn, 128), q))
            }

            function calculateLagrange(pMem) {
                let w := 1

                mstore(
                    add(pMem, pEval_l1),
                    mulmod(n, mod(add(sub(mload(add(pMem, pXi)), w), q), q), q)
                )

                inverseArray(add(pMem, pZhInv), 2)

                let zh := mload(add(pMem, pZh))
                w := 1

                mstore(add(pMem, pEval_l1), mulmod(mload(add(pMem, pEval_l1)), zh, q))
            }

            function calculatePI(pMem, pPub) {
                let pl := 0

                pl := mod(
                    add(
                        sub(pl, mulmod(mload(add(pMem, pEval_l1)), calldataload(add(pPub, 0)), q)),
                        q
                    ),
                    q
                )

                mstore(add(pMem, pPI), pl)
            }

            function calculateR0(pMem) {
                let e1 := mload(add(pMem, pPI))

                let e2 := mulmod(mload(add(pMem, pEval_l1)), mload(add(pMem, pAlpha2)), q)

                let e3a :=
                    addmod(
                        calldataload(pEval_a),
                        mulmod(mload(add(pMem, pBeta)), calldataload(pEval_s1), q),
                        q
                    )
                e3a := addmod(e3a, mload(add(pMem, pGamma)), q)

                let e3b :=
                    addmod(
                        calldataload(pEval_b),
                        mulmod(mload(add(pMem, pBeta)), calldataload(pEval_s2), q),
                        q
                    )
                e3b := addmod(e3b, mload(add(pMem, pGamma)), q)

                let e3c := addmod(calldataload(pEval_c), mload(add(pMem, pGamma)), q)

                let e3 := mulmod(mulmod(e3a, e3b, q), e3c, q)
                e3 := mulmod(e3, calldataload(pEval_zw), q)
                e3 := mulmod(e3, mload(add(pMem, pAlpha)), q)

                let r0 := addmod(e1, mod(sub(q, e2), q), q)
                r0 := addmod(r0, mod(sub(q, e3), q), q)

                mstore(add(pMem, pEval_r0), r0)
            }

            function g1_set(pR, pP) {
                mstore(pR, mload(pP))
                mstore(add(pR, 32), mload(add(pP, 32)))
            }

            function g1_setC(pR, x, y) {
                mstore(pR, x)
                mstore(add(pR, 32), y)
            }

            function g1_calldataSet(pR, pP) {
                mstore(pR, calldataload(pP))
                mstore(add(pR, 32), calldataload(add(pP, 32)))
            }

            function g1_acc(pR, pP) {
                let mIn := mload(0x40)
                mstore(mIn, mload(pR))
                mstore(add(mIn, 32), mload(add(pR, 32)))
                mstore(add(mIn, 64), mload(pP))
                mstore(add(mIn, 96), mload(add(pP, 32)))

                let success := staticcall(sub(gas(), 2000), 6, mIn, 128, pR, 64)

                if iszero(success) {
                    mstore(0, 0)
                    return(0, 0x20)
                }
            }

            function g1_mulAcc(pR, pP, s) {
                let success
                let mIn := mload(0x40)
                mstore(mIn, mload(pP))
                mstore(add(mIn, 32), mload(add(pP, 32)))
                mstore(add(mIn, 64), s)

                success := staticcall(sub(gas(), 2000), 7, mIn, 96, mIn, 64)

                if iszero(success) {
                    mstore(0, 0)
                    return(0, 0x20)
                }

                mstore(add(mIn, 64), mload(pR))
                mstore(add(mIn, 96), mload(add(pR, 32)))

                success := staticcall(sub(gas(), 2000), 6, mIn, 128, pR, 64)

                if iszero(success) {
                    mstore(0, 0)
                    return(0, 0x20)
                }
            }

            function g1_mulAccC(pR, x, y, s) {
                let success
                let mIn := mload(0x40)
                mstore(mIn, x)
                mstore(add(mIn, 32), y)
                mstore(add(mIn, 64), s)

                success := staticcall(sub(gas(), 2000), 7, mIn, 96, mIn, 64)

                if iszero(success) {
                    mstore(0, 0)
                    return(0, 0x20)
                }

                mstore(add(mIn, 64), mload(pR))
                mstore(add(mIn, 96), mload(add(pR, 32)))

                success := staticcall(sub(gas(), 2000), 6, mIn, 128, pR, 64)

                if iszero(success) {
                    mstore(0, 0)
                    return(0, 0x20)
                }
            }

            function g1_mulSetC(pR, x, y, s) {
                let success
                let mIn := mload(0x40)
                mstore(mIn, x)
                mstore(add(mIn, 32), y)
                mstore(add(mIn, 64), s)

                success := staticcall(sub(gas(), 2000), 7, mIn, 96, pR, 64)

                if iszero(success) {
                    mstore(0, 0)
                    return(0, 0x20)
                }
            }

            function g1_mulSet(pR, pP, s) {
                g1_mulSetC(pR, mload(pP), mload(add(pP, 32)), s)
            }

            function calculateD(pMem) {
                let _pD := add(pMem, pD)
                let gamma := mload(add(pMem, pGamma))
                let mIn := mload(0x40)
                mstore(0x40, add(mIn, 256)) // d1, d2, d3 & d4 (4*64 bytes)

                g1_setC(_pD, Qcx, Qcy)
                g1_mulAccC(_pD, Qmx, Qmy, mulmod(calldataload(pEval_a), calldataload(pEval_b), q))
                g1_mulAccC(_pD, Qlx, Qly, calldataload(pEval_a))
                g1_mulAccC(_pD, Qrx, Qry, calldataload(pEval_b))
                g1_mulAccC(_pD, Qox, Qoy, calldataload(pEval_c))

                let betaxi := mload(add(pMem, pBetaXi))
                let val1 := addmod(addmod(calldataload(pEval_a), betaxi, q), gamma, q)

                let val2 :=
                    addmod(addmod(calldataload(pEval_b), mulmod(betaxi, k1, q), q), gamma, q)

                let val3 :=
                    addmod(addmod(calldataload(pEval_c), mulmod(betaxi, k2, q), q), gamma, q)

                let d2a :=
                    mulmod(mulmod(mulmod(val1, val2, q), val3, q), mload(add(pMem, pAlpha)), q)

                let d2b := mulmod(mload(add(pMem, pEval_l1)), mload(add(pMem, pAlpha2)), q)

                // We'll use mIn to save d2
                g1_calldataSet(add(mIn, 192), pZ)
                g1_mulSet(mIn, add(mIn, 192), addmod(addmod(d2a, d2b, q), mload(add(pMem, pU)), q))

                val1 := addmod(
                    addmod(
                        calldataload(pEval_a),
                        mulmod(mload(add(pMem, pBeta)), calldataload(pEval_s1), q),
                        q
                    ),
                    gamma,
                    q
                )

                val2 := addmod(
                    addmod(
                        calldataload(pEval_b),
                        mulmod(mload(add(pMem, pBeta)), calldataload(pEval_s2), q),
                        q
                    ),
                    gamma,
                    q
                )

                val3 := mulmod(
                    mulmod(mload(add(pMem, pAlpha)), mload(add(pMem, pBeta)), q),
                    calldataload(pEval_zw),
                    q
                )

                // We'll use mIn + 64 to save d3
                g1_mulSetC(add(mIn, 64), S3x, S3y, mulmod(mulmod(val1, val2, q), val3, q))

                // We'll use mIn + 128 to save d4
                g1_calldataSet(add(mIn, 128), pT1)

                g1_mulAccC(
                    add(mIn, 128),
                    calldataload(pT2),
                    calldataload(add(pT2, 32)),
                    mload(add(pMem, pXin))
                )
                let xin2 := mulmod(mload(add(pMem, pXin)), mload(add(pMem, pXin)), q)
                g1_mulAccC(add(mIn, 128), calldataload(pT3), calldataload(add(pT3, 32)), xin2)

                g1_mulSetC(
                    add(mIn, 128),
                    mload(add(mIn, 128)),
                    mload(add(mIn, 160)),
                    mload(add(pMem, pZh))
                )

                mstore(add(add(mIn, 64), 32), mod(sub(qf, mload(add(add(mIn, 64), 32))), qf))
                mstore(add(mIn, 160), mod(sub(qf, mload(add(mIn, 160))), qf))
                g1_acc(_pD, mIn)
                g1_acc(_pD, add(mIn, 64))
                g1_acc(_pD, add(mIn, 128))
            }

            function calculateF(pMem) {
                let p := add(pMem, pF)

                g1_set(p, add(pMem, pD))
                g1_mulAccC(p, calldataload(pA), calldataload(add(pA, 32)), mload(add(pMem, pV1)))
                g1_mulAccC(p, calldataload(pB), calldataload(add(pB, 32)), mload(add(pMem, pV2)))
                g1_mulAccC(p, calldataload(pC), calldataload(add(pC, 32)), mload(add(pMem, pV3)))
                g1_mulAccC(p, S1x, S1y, mload(add(pMem, pV4)))
                g1_mulAccC(p, S2x, S2y, mload(add(pMem, pV5)))
            }

            function calculateE(pMem) {
                let s := mod(sub(q, mload(add(pMem, pEval_r0))), q)

                s := addmod(s, mulmod(calldataload(pEval_a), mload(add(pMem, pV1)), q), q)
                s := addmod(s, mulmod(calldataload(pEval_b), mload(add(pMem, pV2)), q), q)
                s := addmod(s, mulmod(calldataload(pEval_c), mload(add(pMem, pV3)), q), q)
                s := addmod(s, mulmod(calldataload(pEval_s1), mload(add(pMem, pV4)), q), q)
                s := addmod(s, mulmod(calldataload(pEval_s2), mload(add(pMem, pV5)), q), q)
                s := addmod(s, mulmod(calldataload(pEval_zw), mload(add(pMem, pU)), q), q)

                g1_mulSetC(add(pMem, pE), G1x, G1y, s)
            }

            function checkPairing(pMem) -> isOk {
                let mIn := mload(0x40)
                mstore(0x40, add(mIn, 576)) // [0..383] = pairing data, [384..447] = pWxi, [448..512] = pWxiw

                let _pWxi := add(mIn, 384)
                let _pWxiw := add(mIn, 448)
                let _aux := add(mIn, 512)

                g1_calldataSet(_pWxi, pWxi)
                g1_calldataSet(_pWxiw, pWxiw)

                // A1
                g1_mulSet(mIn, _pWxiw, mload(add(pMem, pU)))
                g1_acc(mIn, _pWxi)
                mstore(add(mIn, 32), mod(sub(qf, mload(add(mIn, 32))), qf))

                // [X]_2
                mstore(add(mIn, 64), X2x2)
                mstore(add(mIn, 96), X2x1)
                mstore(add(mIn, 128), X2y2)
                mstore(add(mIn, 160), X2y1)

                // B1
                g1_mulSet(add(mIn, 192), _pWxi, mload(add(pMem, pXi)))

                let s := mulmod(mload(add(pMem, pU)), mload(add(pMem, pXi)), q)
                s := mulmod(s, w1, q)
                g1_mulSet(_aux, _pWxiw, s)
                g1_acc(add(mIn, 192), _aux)
                g1_acc(add(mIn, 192), add(pMem, pF))
                mstore(add(pMem, add(pE, 32)), mod(sub(qf, mload(add(pMem, add(pE, 32)))), qf))
                g1_acc(add(mIn, 192), add(pMem, pE))

                // [1]_2
                mstore(add(mIn, 256), G2x2)
                mstore(add(mIn, 288), G2x1)
                mstore(add(mIn, 320), G2y2)
                mstore(add(mIn, 352), G2y1)

                let success := staticcall(sub(gas(), 2000), 8, mIn, 384, mIn, 0x20)

                isOk := and(success, mload(mIn))
            }

            let pMem := mload(0x40)
            mstore(0x40, add(pMem, lastMem))

            checkProofData()
            calculateChallenges(pMem, _pubSignals)
            calculateLagrange(pMem)
            calculatePI(pMem, _pubSignals)
            calculateR0(pMem)
            calculateD(pMem)
            calculateF(pMem)
            calculateE(pMem)
            let isValid := checkPairing(pMem)

            mstore(0x40, sub(pMem, lastMem))
            mstore(0, isValid)
            return(0, 0x20)
        }
    }
}
