'use client';
import { formatGraphLabel, timeOptions } from '@/app/utils/graph';
import { Label } from '@/lib/Label';
import { BarChart } from '@tremor/react';
import { useEffect, useState } from 'react';
import ReactSelect from 'react-select';
import aggregateCountsByInterval from '../../../[mirrorId]/aggregatedCountsByInterval';

type QrepStatusRow = {
  partitionID: string;
  startTime: Date | null;
  endTime: Date | null;
  numRows: number | null;
};

function QrepGraph({ syncs }: { syncs: QrepStatusRow[] }) {
  let [aggregateType, setAggregateType] = useState('hour');
  const initialCount: [string, number][] = [];
  let [counts, setCounts] = useState(initialCount);

  useEffect(() => {
    let rows = syncs.map((sync) => ({
      timestamp: sync.startTime!,
      count: sync.numRows ?? 0,
    }));

    let counts = aggregateCountsByInterval(rows, aggregateType);
    counts = counts.slice(0, 29);
    counts = counts.reverse();
    setCounts(counts);
  }, [aggregateType, syncs]);

  return (
    <div>
      <div className='float-right'>
        <ReactSelect
          id={aggregateType}
          placeholder='Select a timeframe'
          options={timeOptions}
          defaultValue={{ label: 'hour', value: 'hour' }}
          onChange={(val, _) => val && setAggregateType(val.value)}
        />
      </div>
      <div style={{ height: '3rem' }}>
        <Label variant='headline'>Partition sync history</Label>
      </div>
      <BarChart
        className='mt-3'
        data={counts.map((count) => ({
          name: formatGraphLabel(new Date(count[0]), aggregateType),
          'Rows synced at a point in time': count[1],
        }))}
        index='name'
        categories={['Rows synced at a point in time']}
      />
    </div>
  );
}

export default QrepGraph;
