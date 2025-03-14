import pandas as pd
import os

def main():
    sessions_df, queries_df = read_migration_files()

    ip_to_uuid = dict()

    for index, row in sessions_df.iterrows():
        ip = row['ip']
        uuid = row['uuid']
        if ip not in ip_to_uuid:
            ip_to_uuid[ip] = []
        ip_to_uuid[ip].append(uuid)

    print("\nIP to UUID mapping:")
    for ip, uuids in ip_to_uuid.items():
        print(f"{ip}: {uuids}")

def read_migration_files():
    # Get current directory (where migration.py is located)
    migration_dir = os.path.dirname(os.path.abspath(__file__))
    
    # Read the CSV files
    try:
        sessions_df = pd.read_csv(os.path.join(migration_dir, 'test.sessions.0000000010000.csv'))
        queries_df = pd.read_csv(os.path.join(migration_dir, 'test.queries.0000000010000.csv'))
        
        print(f"Successfully loaded {len(sessions_df)} sessions")
        print(f"Successfully loaded {len(queries_df)} queries")
        
        # Preview the data
        print("\nSessions Preview:")
        print(sessions_df.head())
        print("\nQueries Preview:")
        print(queries_df.head())
        
        return sessions_df, queries_df
        
    except FileNotFoundError as e:
        print(f"Error: Could not find CSV file - {e}")
    except pd.errors.EmptyDataError:
        print("Error: One or both CSV files are empty")
    except Exception as e:
        print(f"Error loading CSV files: {e}")

if __name__ == "__main__":
    main()