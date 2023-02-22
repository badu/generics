# Timsort

This package was found [here](https://github.com/psilva261/timsort) and ported to generics.
This implementation was initially ported to Go by [Mike Kroutikov](https://github.com/pgmmpk)
and derived from Java's TimSort object by Josh Bloch,
which, in turn, was based on
the [original code by Tim Peters](https://svn.python.org/projects/python/trunk/Objects/listsort.txt).

# Observations

As Go evolved (the original repository was at Go 1.14) the performance got better, therefor the standard "sort" package
keeps getting better.

Another observation is that sometimes the generic version beats the original one, and other times doesn't. However, the
standard package has average number of bytes allocated per operation and the number of allocations per operation equal
to zero for all the benchmarks (`go test -v -benchmem ./... -bench .`)

## Benchmark results, after porting to generic compared with original
---

goos: linux
goarch: amd64
go version: 1.20
cpu: Intel(R) Core(TM) i7-4770 CPU @ 3.40GHz

#    

| Benchmark                                                   |     Runs |            Time |
|:------------------------------------------------------------|---------:|----------------:|
| [Generics] BenchmarkTimsortStructsXOR100-8                	 |  352920	 |      3559 ns/op |
| [Generics] BenchmarkTimsortInterfacesXOR100-8             	 |  228820	 |      4989 ns/op |
| [--------] BenchmarkStandardStructsXOR100-8               	 |  220815	 |      5657 ns/op |
| [Original] BenchmarkTimsortXor100-8                 	       |  269606	 |      4567 ns/op |
| [Original] BenchmarkTimsortInterXor100-8            	       |  219492	 |      5348 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStructsSorted100-8             	 | 1254387	 |     962.5 ns/op |
| [Generics] BenchmarkTimsortInterfacesSorted100-8          	 |  846442	 |      1434 ns/op |
| [--------] BenchmarkStandardStructsSorted100-8            	 | 1646575	 |     736.9 ns/op |
| [Original] BenchmarkTimsortSorted100-8              	       | 1000000	 |      1214 ns/op |
| [Original] BenchmarkTimsortInterSorted100-8         	       |  749102	 |      1528 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStructsReverseSorted100-8      	 | 1000000	 |      1057 ns/op |
| [Generics] BenchmarkTimsortInterfacesReverseSorted100-8   	 |  705547	 |      1730 ns/op |
| [--------] BenchmarkStandardStructsReverseSorted100-8     	 |  155811	 |      7362 ns/op |
| [Original] BenchmarkTimsortRevSorted100-8           	       |  882352	 |      1456 ns/op |
| [Original] BenchmarkTimsortInterRevSorted100-8      	       |  670225	 |      1844 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStructsRandom100-8             	 |  175089	 |      6755 ns/op |
| [Generics] BenchmarkTimsortInterfacesRandom100-8          	 |  141384	 |      8233 ns/op |
| [--------] BenchmarkStandardStructsRandom100-8            	 |  129858	 |      9550 ns/op |
| [Original] BenchmarkTimsortRandom100-8              	       |  228678	 |      5553 ns/op |
| [Original] BenchmarkTimsortInterRandom100-8         	       |  212698	 |      5746 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStructsXOR1K-8                 	 |   35136	 |     35509 ns/op |
| [Generics] BenchmarkTimsortInterfacesXOR1K-8              	 |   23154	 |     52038 ns/op |
| [--------] BenchmarkStandardStructsXOR1K-8                	 |   10000	 |    104291 ns/op |
| [Original] BenchmarkTimsortXor1K-8                  	       |   26025	 |     49134 ns/op |
| [Original] BenchmarkTimsortInterXor1K-8             	       |   22719	 |     55499 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStructsSorted1K-8              	 |  258412	 |      4623 ns/op |
| [Generics] BenchmarkTimsortInterfacesSorted1K-8           	 |  133430	 |      8458 ns/op |
| [--------] BenchmarkStandardStructsSorted1K-8             	 |  216129	 |      5740 ns/op |
| [Original] BenchmarkTimsortSorted1K-8               	       |  208828	 |      6120 ns/op |
| [Original] BenchmarkTimsortInterSorted1K-8          	       |  136764	 |      8564 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStructsReverseSorted1K-8       	 |  222482	 |      5386 ns/op |
| [Generics] BenchmarkTimsortInterfacesReverseSorted1K-8    	 |  106596	 |     10538 ns/op |
| [--------] BenchmarkStandardStructsReverseSorted1K-8      	 |   13818	 |     86140 ns/op |
| [Original] BenchmarkTimsortRevSorted1K-8            	       |  160087	 |      8089 ns/op |
| [Original] BenchmarkTimsortInterRevSorted1K-8       	       |  107556	 |     10947 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStructsRandom1K-8              	 |   10000	 |    102048 ns/op |
| [Generics] BenchmarkTimsortInterfacesRandom1K-8           	 |   10000	 |    123271 ns/op |
| [--------] BenchmarkStandardStructsRandom1K-8             	 |    6454	 |    191541 ns/op |
| [Original] BenchmarkTimsortRandom1K-8               	       |   10000	 |    116834 ns/op |
| [Original] BenchmarkTimsortInterRandom1K-8          	       |   10000	 |    118543 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStructsXOR1M-8                 	 |      18	 |  65165792 ns/op |
| [Generics] BenchmarkTimsortInterfacesXOR1M-8              	 |      10	 | 103312919 ns/op |
| [--------] BenchmarkStandardStructsXOR1M-8                	 |       6	 | 216099678 ns/op |
| [Original] BenchmarkTimsortXor1M-8                  	       |      12	 |  99730054 ns/op |
| [Original] BenchmarkTimsortInterXor1M-8             	       |      10	 | 129135147 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStructsSorted1M-8              	 |     469	 |   2659028 ns/op |
| [Generics] BenchmarkTimsortInterfacesSorted1M-8           	 |     141	 |   8581116 ns/op |
| [--------] BenchmarkStandardStructsSorted1M-8             	 |     148	 |   8058755 ns/op |
| [Original] BenchmarkTimsortSorted1M-8               	       |     272	 |   4430028 ns/op |
| [Original] BenchmarkTimsortInterSorted1M-8          	       |     151	 |   8038767 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStructsReverseSorted1M-8       	 |     100	 |  10725008 ns/op |
| [Generics] BenchmarkTimsortInterfacesRevSorted1M-8        	 |     100	 |  11316071 ns/op |
| [--------] BenchmarkStandardStructsReverseSorted1M-8      	 |      10	 | 100150976 ns/op |
| [Original] BenchmarkTimsortInterRevSorted1M-8       	       |     100	 |  10216143 ns/op |
| [Original] BenchmarkTimsortRevSorted1M-8            	       |     194	 |   6443761 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStructsRandom1M-8              	 |       4	 | 259547227 ns/op |
| [Generics] BenchmarkTimsortInterfacesRandom1M-8           	 |       3	 | 369764016 ns/op |
| [--------] BenchmarkStandardStructsRandom1M-8             	 |       2	 | 567693604 ns/op |
| [Original] BenchmarkTimsortInterRandom1M-8          	       |       3	 | 387636301 ns/op |
| [Original] BenchmarkTimsortRandom1M-8               	       |       3	 | 358330702 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortIntsXOR100-8                   	 |  329522	 |      3221 ns/op |
| [--------] BenchmarkStandardIntsXOR100-8                  	 |  379278	 |      2912 ns/op |
| [Original] BenchmarkTimsortIXor100-8                	       |  337268	 |      3556 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortIntsSorted100-8                	 | 1561195	 |     768.3 ns/op |
| [--------] BenchmarkStandardIntsSorted100-8               	 | 2400471	 |     501.8 ns/op |
| [Original] BenchmarkTimsortISorted100-8             	       | 1406851	 |     843.4 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortIntsReverseSorted100-8         	 | 1344094	 |     931.1 ns/op |
| [--------] BenchmarkStandardIntsReverseSorted100-8        	 | 1779391	 |     683.9 ns/op |
| [Original] BenchmarkTimsortIRevSorted100-8          	       | 1220152	 |     971.2 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortIntsRandom100-8                	 |  180652	 |      6100 ns/op |
| [--------] BenchmarkStandardIntsRandom100-8               	 |  215168	 |      6483 ns/op |
| [Original] BenchmarkTimsortIRandom100-8             	       |  309997	 |      3506 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortIntsXOR1K-8                    	 |   32005	 |     38292 ns/op |
| [--------] BenchmarkStandardIntsXOR1K-8                   	 |   18420	 |     63981 ns/op |
| [Original] BenchmarkTimsortIXor1K-8                 	       |   39088	 |     35184 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortIntsSorted1K-8                 	 |  315118	 |      3726 ns/op |
| [--------] BenchmarkStandardIntsSorted1K-8                	 |  282036	 |      3811 ns/op |
| [Original] BenchmarkTimsortISorted1K-8              	       |  324958	 |      3485 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortIntsReverseSorted1K-8          	 |  289677	 |      4183 ns/op |
| [--------] BenchmarkStandardIntsReverseSorted1K-8         	 |  243510	 |      4919 ns/op |
| [Original] BenchmarkTimsortIRevSorted1K-8           	       |  289059	 |      4338 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortIntsRandom1K-8                 	 |   13971	 |     84148 ns/op |
| [--------] BenchmarkStandardIntsRandom1K-8                	 |   14886	 |     82114 ns/op |
| [Original] BenchmarkTimsortIRandom1K-8              	       |   15381	 |     78625 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortIntsXOR1M-8                    	 |      21	 |  53455832 ns/op |
| [--------] BenchmarkStandardIntsXOR1M-8                   	 |      21	 |  50006838 ns/op |
| [Original] BenchmarkTimsortIXor1M-8                 	       |      18	 |  55598534 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortIntsSorted1M-8                 	 |     379	 |   2690324 ns/op |
| [--------] BenchmarkStandardIntsSorted1M-8                	 |     373	 |   3522509 ns/op |
| [Original] BenchmarkTimsortISorted1M-8              	       |     370	 |   2820619 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortIntsReverseSorted1M-8          	 |     318	 |   3749591 ns/op |
| [--------] BenchmarkStandardIntsReverseSorted1M-8         	 |     213	 |   5387459 ns/op |
| [Original] BenchmarkTimsortIRevSorted1M-8           	       |     334	 |   3510574 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortIntsRandom1M-8                 	 |       7	 | 160864569 ns/op |
| [--------] BenchmarkStandardIntsRandom1M-8                	 |       7	 | 152528051 ns/op |
| [Original] BenchmarkTimsortIRandom1M-8              	       |       7	 | 166340788 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStringsXOR100-8                	 |  173881	 |      7264 ns/op |
| [--------] BenchmarkStandardStringsXOR100-8               	 |  176869	 |      6591 ns/op |
| [Original] BenchmarkTimsortStrXor100-8              	       |  161359	 |      7171 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStringsSorted100-8             	 |  230437	 |      4656 ns/op |
| [--------] BenchmarkStandardStringsSorted100-8            	 |  192972	 |      6696 ns/op |
| [Original] BenchmarkTimsortStrSorted100-8           	       |  246375	 |      4531 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStringsReverseSorted100-8      	 |  227742	 |      5476 ns/op |
| [--------] BenchmarkStandardStringsReversedSorted100-8    	 |  124449	 |     10107 ns/op |
| [Original] BenchmarkTimsortStrRevSorted100-8        	       |  237786	 |      5360 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStringsRandom100-8             	 |  100462	 |     12019 ns/op |
| [--------] BenchmarkStandardStringsRandom100-8            	 |  118444	 |      9828 ns/op |
| [Original] BenchmarkTimsortStrRandom100-8           	       |  149917	 |      8273 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStringsXOR1K-8                 	 |   12228	 |    106983 ns/op |
| [--------] BenchmarkStandardStringsXOR1K-8                	 |    7848	 |    139966 ns/op |
| [Original] BenchmarkTimsortStrXor1K-8               	       |   12004	 |     96417 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStringsSorted1K-8              	 |   33210	 |     32549 ns/op |
| [--------] BenchmarkStandardStringsSorted1K-8             	 |    7804	 |    142907 ns/op |
| [Original] BenchmarkTimsortStrSorted1K-8            	       |   39640	 |     30287 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStringsReverseSorted1K-8       	 |   35281	 |     35295 ns/op |
| [--------] BenchmarkStandardStringsReverseSorted1K-8      	 |    8728	 |    147519 ns/op |
| [Original] BenchmarkTimsortStrRevSorted1K-8         	       |   34330	 |     35371 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStringsRandom1K-8              	 |    7078	 |    181304 ns/op |
| [--------] BenchmarkStandardStringsRandom1K-8             	 |    8503	 |    145830 ns/op |
| [Original] BenchmarkTimsortStrRandom1K-8            	       |    7066	 |    184913 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStringsXOR1M-8                 	 |       5	 | 211638297 ns/op |
| [--------] BenchmarkStandardStringsXOR1M-8                	 |      10	 | 113652126 ns/op |
| [Original] BenchmarkTimsortStrXor1M-8               	       |       5	 | 219005952 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStringsSorted1M-8              	 |      10	 | 113723975 ns/op |
| [--------] BenchmarkStandardStringsSorted1M-8             	 |       3	 | 344755139 ns/op |
| [Original] BenchmarkTimsortStrSorted1M-8            	       |      10	 | 102599597 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStringsReverseSorted1M-8       	 |      10	 | 110395080 ns/op |
| [--------] BenchmarkStandardStringsReverseSorted1M-8      	 |       5	 | 250614559 ns/op |
| [Original] BenchmarkTimsortStrRevSorted1M-8         	       |      10	 | 102154311 ns/op |
|                                                             |          |                 |
| [Generics] BenchmarkTimsortStringsRandom1M-8              	 |       2	 | 785339384 ns/op |
| [--------] BenchmarkStandardStringsRandom1M-8             	 |       3	 | 374392649 ns/op |
| [Original] BenchmarkTimsortStrRandom1M-8            	       |       2	 | 725358523 ns/op |

## Examples, adapted to generics

### As drop-in replacement for sort.Sort

    package main

    import (
		"github.com/badu/generics/timsort"
		"fmt"
		"sort"
    )

    func main() {
		l := []string{"c", "a", "b"}
		timsort.TimSort(sort.StringSlice(l)
		fmt.Printf("sorted array: %+v\n", l)
    }

### Explicit "less" function

	package main

	import (
		"github.com/badu/generics/timsort"
		"fmt"
	)

	type Record struct {
		ssn  int
		name string
	}

	func BySsn(a, b Record) bool {
		return a.ssn < b.ssn
	}

	func ByName(a, b Record) bool {
		return a.name < b.name
	}

	func main() {
		db := make([]Record, 3)
		db[0] = Record{123456789, "joe"}
		db[1] = Record{101765430, "sue"}
		db[2] = Record{345623452, "mary"}

		// sorts array by ssn (ascending)
		timsort.Sort(db, BySsn)
		fmt.Printf("sorted by ssn: %v\n", db)

		// now re-sort same array by name (ascending)
		timsort.Sort(db, ByName)
		fmt.Printf("sorted by name: %v\n", db)
	}
