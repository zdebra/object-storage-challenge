package storagegateway

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestFindInstanceForKey(t *testing.T) {
	service := NewService(&StorageInstance{id: "test1"}, &StorageInstance{id: "test2"}, &StorageInstance{id: "test3"})

	t.Run("same instance for same key", func(t *testing.T) {
		firstInstance, err := service.findInstanceForKey("x")
		require.NoError(t, err)
		secondInstance, err := service.findInstanceForKey("x")
		require.NoError(t, err)

		require.Equal(t, firstInstance, secondInstance)
	})

	t.Run("distribution", func(t *testing.T) {
		keys := []string{
			"6808683896",
			"1722217355",
			"2014100978",
			"2492790589",
			"7529035987",
			"4441906119",
			"8220669715",
			"8450683636",
			"6503458644",
			"3740900001",
			"1147860863",
			"4066345700",
			"5184022855",
			"4823061950",
			"7477368853",
			"8312693379",
			"8051931031",
			"8082454144",
			"8731991045",
			"5459906428",
			"0180379073",
			"7483559881",
			"7034490252",
			"7608771682",
			"3166593755",
			"7151201567",
			"9114859262",
			"0813118779",
			"7260860460",
			"4901323337",
			"1560589010",
			"7607129374",
			"1902733421",
			"6183620571",
			"1941718405",
			"1694022399",
			"4887283230",
			"4065801220",
			"2296087803",
			"0796386003",
			"4200045484",
			"5754406879",
			"1323336396",
			"4166783382",
			"1476950248",
			"9067219971",
			"7371440305",
			"7838285875",
			"0328582062",
			"2700825384",
			"2787550956",
			"9395568883",
			"7414067849",
			"2510171543",
			"8359959477",
			"0401031077",
			"9525144785",
			"3626993131",
			"9343011541",
			"4931420026",
			"1219956059",
			"0451410826",
			"3712270881",
			"5934658431",
			"0117986722",
			"3938994955",
			"6883098396",
			"9850335621",
			"4201722528",
			"5266893839",
			"0937869344",
			"7091474805",
			"3249055632",
			"5115543810",
			"4638008620",
			"7068370171",
			"2853122240",
			"3640722021",
			"5421680398",
			"6313634510",
			"5226756134",
			"5155989654",
			"4878653796",
			"5355524157",
			"6939244921",
			"4194420594",
			"7240101540",
			"0081522779",
			"8856844778",
			"3844612814",
			"6851498261",
			"4340785259",
			"6738792703",
			"9747591423",
			"0979191103",
			"5457586165",
			"0790078932",
			"4191371059",
			"2720972906",
			"4273297698",
		}

		instancesCnt := map[string]int{}
		for _, key := range keys {
			instance, err := service.findInstanceForKey(key)
			require.NoError(t, err)
			require.NotNil(t, instance)
			instancesCnt[instance.id]++
		}

		// this is not very good test as it depends on the hash function
		// but I wanted to understand how the distribution looks like
		assert.Equal(t, 28, instancesCnt["test1"])
		assert.Equal(t, 47, instancesCnt["test2"])
		assert.Equal(t, 25, instancesCnt["test3"])
	})
}
